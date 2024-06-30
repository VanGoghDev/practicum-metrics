package pgstorage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"
)

type PgStorage struct {
	zlog *zap.SugaredLogger
	pool *pgxpool.Pool
}

func New(ctx context.Context, zlog *zap.SugaredLogger, cfg *config.Config) (*PgStorage, error) {
	pool, err := pgxpool.New(ctx, cfg.DBConnectionString)
	if err != nil {
		zlog.Warnf("failed to establich connection with db: %w:", err)
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	s := &PgStorage{
		zlog: zlog,
		pool: pool,
	}

	err = s.runMigrations(cfg.DBConnectionString)
	if err != nil {
		zlog.Warnf("failed to create schema: %w", err)
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	err = s.pingWithTimeout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping with timeout: %w", err)
	}
	return s, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func (s *PgStorage) runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		d,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to get a create new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func (s *PgStorage) SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				s.zlog.Errorf("failed to rollaback the transaction: %w", err)
			}
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
