package filesaver

import (
	"fmt"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage"
	"go.uber.org/zap"
)

type FileStorage interface {
	Save(metrics []*models.Metrics) error
	GetMetrics() ([]*models.Metrics, error)
}

type MemStorage interface {
	GetMetrics() ([]*models.Metrics, error)
	SaveMetrics([]*models.Metrics) error
}

type FileSaver struct {
	logger        *zap.Logger
	filestrg      FileStorage
	memstorage    MemStorage
	storeInterval time.Duration
	restore       bool
}

func New(cfg *config.Config, log *zap.Logger, mstrg MemStorage) (*FileSaver, error) {
	filestrg, err := filestorage.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file storage: %w", err)
	}
	return &FileSaver{
		logger:        log,
		memstorage:    mstrg,
		filestrg:      filestrg,
		storeInterval: cfg.StoreInterval,
		restore:       cfg.Restore,
	}, nil
}

func (fs FileSaver) Run() {
	if fs.restore {
		err := fs.Restore()
		if err != nil {
			fs.logger.Warn(fmt.Sprintf("failed to restore metrics: %v", err))
		}
	}

	storeTicker := time.NewTicker(fs.storeInterval)
	defer storeTicker.Stop()
	quit := make(chan struct{})

	for {
		select {
		case <-storeTicker.C:
			metrics, err := fs.memstorage.GetMetrics()
			if err != nil {
				fs.logger.Warn(fmt.Sprintf("filestorage returned error on get: %v", err))
			}

			err = fs.filestrg.Save(metrics)
			if err != nil {
				fs.logger.Warn(fmt.Sprintf("filestorage returned error on write: %v", err))
			}
		case <-quit:
			storeTicker.Stop()
			return
		}
	}
}

func (fs FileSaver) Restore() error {
	metrics, err := fs.filestrg.GetMetrics()
	if err != nil {
		return fmt.Errorf("failed to restore: %w", err)
	}
	err = fs.memstorage.SaveMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to save metrics to memory: %w", err)
	}
	return nil
}
