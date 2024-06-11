package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

type PgStorage struct {
	db *sql.DB
}

const (
	stmtErrMsg = "failed to close stmt"
	pingErrMsg = "failed to connect to db with timeout"
)

func New(ctx context.Context, zlog *zap.SugaredLogger, cfg *config.Config) (*PgStorage, error) {
	db, err := sql.Open("pgx", cfg.DBConnectionString)
	if err != nil {
		zlog.Warnf("failed to open db: %v", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	err = pingWithTimeout(db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = createSchema(ctx, db)
	if err != nil {
		zlog.Warnf("failed to create schema: %w", err)
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &PgStorage{
		db: db,
	}, nil
}

func (s *PgStorage) SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error) {
	err = pingWithTimeout(s.db)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}

	deleteStmt, err := tx.PrepareContext(ctx, "DELETE FROM metrics WHERE name = $1")
	defer func() {
		err = deleteStmt.Close()
	}()
	if err != nil {
		return fmt.Errorf("failed to init delete statement: %w", err)
	}

	insrtStmt, err := tx.PrepareContext(ctx, "INSERT INTO metrics (name, g_type, g_value, delta) VALUES ($1, $2, $3, $4)")
	defer func() {
		err = insrtStmt.Close()
	}()
	if err != nil {
		return fmt.Errorf("failed to init statement: %w", err)
	}

	updStmt, err := tx.PrepareContext(ctx, "INSERT INTO metrics(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)"+
		" ON CONFLICT(name) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta")
	defer func() {
		err = updStmt.Close()
	}()
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
			_, err := updStmt.ExecContext(ctx, v.ID, v.MType, defaultValue, defaultDelta)
			if err != nil {
				err := tx.Rollback()
				return fmt.Errorf("failed to execute insert statement: %w", err)
			}
		case handlers.Gauge:
			_, err := deleteStmt.ExecContext(ctx, v.ID)
			if err != nil {
				err := tx.Rollback()
				return fmt.Errorf("failed to execute delete statement: %w", err)
			}

			_, err = insrtStmt.ExecContext(ctx, v.ID, v.MType, defaultValue, defaultDelta)
			if err != nil {
				err := tx.Rollback()
				return fmt.Errorf("failed to execute insert statement: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PgStorage) SaveGauge(ctx context.Context, name string, value float64) (err error) {
	err = pingWithTimeout(s.db)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	stmt, err := s.db.Prepare("INSERT INTO metrics(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)")
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("%s: %w", stmtErrMsg, err)
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to prepareContext: %w", err)
	}

	_, err = stmt.ExecContext(ctx, name, handlers.Gauge, value, 0)
	if err != nil {
		return fmt.Errorf("failed to execute save querry: %w", err)
	}
	return nil
}

func (s *PgStorage) SaveCount(ctx context.Context, name string, value int64) (err error) {
	err = pingWithTimeout(s.db)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO metrics(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)"+
		" ON CONFLICT(name) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta")
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("%s: %w", stmtErrMsg, err)
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to prepareContext: %w", err)
	}

	_, err = stmt.ExecContext(ctx, name, handlers.Counter, 0, value)
	if err != nil {
		return fmt.Errorf("failed to execute save querry: %w", err)
	}
	return nil
}

func (s *PgStorage) Gauges(ctx context.Context) (gauges []models.Gauge, err error) {
	err = pingWithTimeout(s.db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stmt, err := s.db.Prepare("SELECT name, g_value FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare gauges query: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("%s: %w", stmtErrMsg, err)
		}
	}()

	rows, err := stmt.QueryContext(ctx)
	defer func() {
		err = rows.Close()
		if err != nil {
			err = fmt.Errorf("failed to close row: %w", err)
		}
	}()

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
	err = pingWithTimeout(s.db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	stmt, err := s.db.Prepare("SELECT name, delta FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare counters query: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("failed to close stmt: %w", err)
		}
	}()

	rows, err := stmt.QueryContext(ctx)
	defer func() {
		err = rows.Close()
		if err != nil {
			err = fmt.Errorf("failed to close row: %w", err)
		}
	}()

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
	err = pingWithTimeout(s.db)
	if err != nil {
		return models.Gauge{}, fmt.Errorf("%w", err)
	}

	stmt, err := s.db.Prepare("SELECT name, g_value FROM metrics WHERE name = $1")
	if err != nil {
		return models.Gauge{}, fmt.Errorf("failed to prepare counter query: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("failed to close stmt: %w", err)
		}
	}()

	row := stmt.QueryRowContext(ctx, name)
	err = row.Scan(&gauge.Name, &gauge.Value)
	if err != nil {
		return models.Gauge{}, fmt.Errorf("failed to scan counter: %w", err)
	}
	return gauge, nil
}

func (s *PgStorage) Counter(ctx context.Context, name string) (counter models.Counter, err error) {
	err = pingWithTimeout(s.db)
	if err != nil {
		return models.Counter{}, fmt.Errorf("%w", err)
	}

	stmt, err := s.db.Prepare("SELECT name, delta FROM metrics WHERE name = $1")
	if err != nil {
		return models.Counter{}, fmt.Errorf("failed to prepare counter query: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			err = fmt.Errorf("failed to close stmt: %w", err)
		}
	}()

	row := stmt.QueryRowContext(ctx, name)
	err = row.Scan(&counter.Name, &counter.Value)
	if err != nil {
		return models.Counter{}, fmt.Errorf("failed to scan counter: %w", err)
	}
	return counter, nil
}

// Пингует базу. Если что-то не так, то будет пинговать три раза, затем вернет ошибку если недопингуется.
func pingWithTimeout(db *sql.DB) error {
	retriesCount := 0
	maxRetriesCount := 3
	f := 2
	err := db.Ping()
	if err == nil {
		return nil
	}
	for {
		if retriesCount >= maxRetriesCount {
			err = errors.New("failed to connect to db")
			break
		}
		retriesCount++
		err = db.Ping()
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

func (s *PgStorage) Close(ctx context.Context) error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}

	return nil
}

func createSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS metrics(
			name 	VARCHAR(200) PRIMARY KEY,
			g_type 	VARCHAR(200) NOT NULL,
			g_value DOUBLE PRECISION NOT NULL,
			delta 	bigint NOT NULL,
			UNIQUE(name)
		)`,
	}

	for _, stmt := range createSchemaStmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return err
}
