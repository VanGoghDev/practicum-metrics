package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"
)

type PgStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, zlog *zap.SugaredLogger, cfg *config.Config) (*PgStorage, error) {
	pool, err := pgxpool.New(ctx, cfg.DBConnectionString)
	if err != nil {
		zlog.Warnf("failed to establich connection with db: %w:", err)
	}

	err = createSchema(ctx, pool)
	if err != nil {
		zlog.Warnf("failed to create schema: %w", err)
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	s := &PgStorage{
		pool: pool,
	}

	err = s.pingWithTimeout(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return s, nil
}

func (s *PgStorage) SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("failed to rollback the transaction: %v", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}

	_, err = tx.Prepare(ctx, "delete", "DELETE FROM metrics WHERE name = $1")
	if err != nil {
		return fmt.Errorf("failed to init delete statement: %w", err)
	}

	_, err = tx.Prepare(ctx, "insrtStmt", "INSERT INTO metrics (name, g_type, g_value, delta)"+
		"VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to init statement: %w", err)
	}

	_, err = tx.Prepare(ctx, "updStmt", "INSERT INTO metrics(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)"+
		" ON CONFLICT(name) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta")
	if err != nil {
		return fmt.Errorf("failed to init update statement: %w", err)
	}

	for _, v := range metrics {
		var defaultValue float64 = 0
		var defaultDelta int64 = 0
		if v.Value != nil {
			defaultValue = *v.Value
		}
		if v.Delta != nil {
			defaultDelta = *v.Delta
		}
		switch v.MType {
		case handlers.Counter:
			_, err := tx.Exec(ctx, "updStmt", v.ID, v.MType, defaultValue, defaultDelta)
			if err != nil {
				return fmt.Errorf("failed to execute insert statement: %w", err)
			}
		case handlers.Gauge:
			_, err := tx.Exec(ctx, "delete", v.ID)
			if err != nil {
				return fmt.Errorf("failed to execute delete statement: %w", err)
			}

			_, err = tx.Exec(ctx, "insrtStmt", v.ID, v.MType, defaultValue, defaultDelta)
			if err != nil {
				return fmt.Errorf("failed to execute insert statement: %w", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PgStorage) SaveGauge(ctx context.Context, name string, value float64) (err error) {
	_, err = s.pool.Exec(ctx, "INSERT INTO metrics(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)",
		name, handlers.Gauge, value, 0)
	if err != nil {
		return fmt.Errorf("failed to execute save querry: %w", err)
	}
	return nil
}

func (s *PgStorage) SaveCount(ctx context.Context, name string, value int64) (err error) {
	_, err = s.pool.Exec(ctx, "INSERT INTO metrics(name, g_type, g_value, delta)"+
		"VALUES($1, $2, $3, $4) ON CONFLICT(name) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta",
		name, handlers.Counter, 0, value)
	if err != nil {
		return fmt.Errorf("failed to execute save querry: %w", err)
	}
	return nil
}

func (s *PgStorage) Gauges(ctx context.Context) (gauges []models.Gauge, err error) {
	rows, err := s.pool.Query(ctx, "SELECT name, g_value FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to query gauges: %w", err)
	}

	for rows.Next() {
		var g models.Gauge
		err = rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in rows: %w", err)
		}
		gauges = append(gauges, g)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate through rows: %w", err)
	}
	return gauges, nil
}

func (s *PgStorage) Counters(ctx context.Context) (counters []models.Counter, err error) {
	rows, err := s.pool.Query(ctx, "SELECT name, delta FROM metrics")

	if err != nil {
		return nil, fmt.Errorf("failed to query gauges: %w", err)
	}

	for rows.Next() {
		var g models.Counter
		err = rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in rows: %w", err)
		}
		counters = append(counters, g)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate through rows: %w", err)
	}
	return counters, nil
}

func (s *PgStorage) Gauge(ctx context.Context, name string) (gauge models.Gauge, err error) {
	row := s.pool.QueryRow(ctx, "SELECT name, g_value FROM metrics WHERE name = $1", name)
	err = row.Scan(&gauge.Name, &gauge.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Gauge{}, serrors.ErrNotFound
		}
		return models.Gauge{}, fmt.Errorf("failed to scan gauge: %w", err)
	}
	return gauge, nil
}

func (s *PgStorage) Counter(ctx context.Context, name string) (counter models.Counter, err error) {
	row := s.pool.QueryRow(ctx, "SELECT name, delta FROM metrics WHERE name = $1", name)
	err = row.Scan(&counter.Name, &counter.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Counter{}, serrors.ErrNotFound
		}
		return models.Counter{}, fmt.Errorf("failed to scan counter: %w", err)
	}
	return counter, nil
}

func (s *PgStorage) Ping(ctx context.Context) error {
	err := s.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	return nil
}

func (s *PgStorage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}

// Пингует базу. Если что-то не так, то будет пинговать три раза, затем вернет ошибку если недопингуется.
func (s *PgStorage) pingWithTimeout(ctx context.Context) error {
	retriesCount := 0
	maxRetriesCount := 3
	f := 2
	err := s.pool.Ping(ctx)
	if err == nil {
		return nil
	}
	for {
		if retriesCount >= maxRetriesCount {
			err = errors.New("failed to connect to db")
			break
		}
		retriesCount++
		err = s.Ping(ctx)
		if err == nil {
			break
		}
		interval := (retriesCount * f) - 1
		time.Sleep(time.Duration(interval) * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	return nil
}

func createSchema(ctx context.Context, db *pgxpool.Pool) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("failed to rollback the transaction: %v", err)
		}
	}()

	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS metrics(
			name 	VARCHAR(200) PRIMARY KEY,
			g_type 	VARCHAR(200) NOT NULL,
			g_value DOUBLE PRECISION NOT NULL,
			delta 	bigint NOT NULL,
			UNIQUE(name, g_type)
		)`,
	}

	for _, stmt := range createSchemaStmts {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return err
}
