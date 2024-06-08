package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
)

func New(ctx context.Context, zlog *zap.SugaredLogger, cfg *config.Config) (*PgStorage, error) {
	db, err := sql.Open("pgx", cfg.DBConnectionString)
	if err != nil {
		zlog.Warnf("failed to open db: %v", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		zlog.Warnf("failed to ping db: %v", err)
		return nil, fmt.Errorf("failed to ping db: %w", err)
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO gauges (name, g_type, g_value, delta) VALUES ($1, $2, $3, $4)")
	defer func() {
		err = stmt.Close()
	}()

	if err != nil {
		return fmt.Errorf("failed to init statement: %w", err)
	}

	for _, v := range metrics {
		_, err := stmt.ExecContext(ctx, v.ID, v.MType, v.Delta, v.Value)
		if err != nil {
			err := tx.Rollback()
			return fmt.Errorf("failed to execute insert statement: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PgStorage) SaveGauge(ctx context.Context, name string, value float64) (err error) {
	stmt, err := s.db.Prepare("INSERT INTO gauges(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)")
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
	stmt, err := s.db.Prepare("INSERT INTO gauges(name, g_type, g_value, delta) VALUES($1, $2, $3, $4)")
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
	stmt, err := s.db.Prepare("SELECT name, g_value FROM gauges")
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
		return nil, fmt.Errorf("faile to query gauges: %w", err)
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
	stmt, err := s.db.Prepare("SELECT name, delta FROM gauges")
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
		return nil, fmt.Errorf("faile to query gauges: %w", err)
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
	stmt, err := s.db.Prepare("SELECT name, g_value FROM gauges WHERE name = $1")
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
	stmt, err := s.db.Prepare("SELECT name, g_type, g_value, delta FROM gauges WHERE name = $1")
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

func (s *PgStorage) Close(ctx context.Context) error {
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
		`CREATE TABLE IF NOT EXISTS gauges(
			id 		INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name 	VARCHAR(200),
			g_type 	VARCHAR(200) NOT NULL,
			g_value DOUBLE PRECISION,
			delta 	DOUBLE PRECISION
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
