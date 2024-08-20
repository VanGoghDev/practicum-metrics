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

// Storage интерфейс хранилища метрик.
type Storage interface {

	// SaveMetrics сохраняет метрики.
	SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error)

	// SaveGauge сохраняет метрики типа Gauge.
	SaveGauge(ctx context.Context, name string, value float64) (err error)

	// SaveCount сохраняет метрики типа Count.
	SaveCount(ctx context.Context, name string, value int64) (err error)

	// Gauges возвращает список метрик типа Gauge.
	Gauges(ctx context.Context) (gauges []models.Gauge, err error)

	// Counters возвращает список метрик типа Counters.
	Counters(ctx context.Context) (counters []models.Counter, err error)

	// Gauge возвращает метрику типа Gauge.
	Gauge(ctx context.Context, name string) (gauge models.Gauge, err error)

	// Counter возвращает метрику типа Counter.
	Counter(ctx context.Context, name string) (counter models.Counter, err error)

	// Ping пинг хранилища.
	Ping(ctx context.Context) error

	// Close закрывает соединение с хранилищем.
	Close(ctx context.Context) error
}

// New возвращает новый экземпляр хранилища.
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
