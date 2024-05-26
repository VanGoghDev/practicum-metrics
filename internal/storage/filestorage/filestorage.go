package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"go.uber.org/zap"
)

type FileStorage struct {
	zlog    *zap.Logger
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func New(log *zap.Logger, cfg *config.Config) (*FileStorage, error) {
	var perm fs.FileMode = 0o666
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to  init file storage: %w", err)
	}

	return &FileStorage{
		zlog:    log,
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (f *FileStorage) Save(metrics []*models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}
	err := f.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	for _, v := range metrics {
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal metrics: %w", err)
		}

		_, err = f.writer.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write metrics to file: %w", err)
		}

		err = f.writer.WriteByte('\n')
		if err != nil {
			return fmt.Errorf("failed to write delimiter to file: %w", err)
		}
	}
	err = f.writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}
	return nil
}

func (f *FileStorage) GetMetrics() ([]*models.Metrics, error) {
	f.zlog.Info("getting metrics from file!")
	metrics := make([]*models.Metrics, 0)
	for f.scanner.Scan() {
		metric := &models.Metrics{}
		data := f.scanner.Bytes()
		if len(data) > 0 {
			err := json.Unmarshal(data, metric)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal metric: %w", err)
			}
			metrics = append(metrics, metric)
		}
	}
	f.zlog.Info("got metrics from file!")

	return metrics, nil
}

func (f *FileStorage) Close() error {
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("Filestorage.Close: %w", err)
	}
	return nil
}
