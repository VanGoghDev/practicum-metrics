package storage

import (
	"fmt"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"go.uber.org/zap"
)

type Storage interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
	Gauges() (gauges []models.Gauge, err error)
	Counters() (counters []models.Counter, err error)
	Gauge(name string) (gauge models.Gauge, err error)
	Counter(name string) (counter models.Counter, err error)
	Close() error
}

func New(cfg *config.Config, zlog *zap.Logger) (Storage, error) {
	var s Storage
	if cfg.FileStoragePath == "" {
		s, err := memstorage.New(zlog)
		if err != nil {
			return nil, fmt.Errorf("failed to init memory storage: %w", err)
		}
		return s, nil
	}
	s, err := filestorage.New(zlog, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init file storage: %w", err)
	}
	return s, nil
}
