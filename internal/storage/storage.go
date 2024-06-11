package storage

import (
	"context"
	"fmt"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/pgstorage"
	"go.uber.org/zap"
)

type Storage interface {
	SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error)
	SaveGauge(ctx context.Context, name string, value float64) (err error)
	SaveCount(ctx context.Context, name string, value int64) (err error)
	Gauges(ctx context.Context) (gauges []models.Gauge, err error)
	Counters(ctx context.Context) (counters []models.Counter, err error)
	Gauge(ctx context.Context, name string) (gauge models.Gauge, err error)
	Counter(ctx context.Context, name string) (counter models.Counter, err error)
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(ctx context.Context, cfg *config.Config, zlog *zap.Logger) (Storage, error) {
	var s Storage
	if cfg.DBConnectionString != "" {
		zlog.Debug("Init db storage")
		s, err := pgstorage.New(ctx, zlog.Sugar(), cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init db storage: %w", err)
		}
		return s, nil
	}
	if cfg.FileStoragePath == "" {
		zlog.Debug("Init memory storage")
		s, err := memstorage.New(zlog)
		if err != nil {
			return nil, fmt.Errorf("failed to init memory storage: %w", err)
		}
		return s, nil
	}
	zlog.Debug("Init file storage")
	s, err := filestorage.New(ctx, zlog, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init file storage: %w", err)
	}
	return s, nil
}
