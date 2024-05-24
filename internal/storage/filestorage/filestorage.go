package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func New(cfg *config.Config) (*FileStorage, error) {
	var mode fs.FileMode = 0o600
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to  init file storage: %w", err)
	}

	return &FileStorage{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (f *FileStorage) Save(metrics []*models.Metrics) error {
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
	err := f.writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}
	return nil
}

func (f *FileStorage) GetMetrics() ([]*models.Metrics, error) {
	metrics := make([]*models.Metrics, 0)
	for f.scanner.Scan() {
		metric := &models.Metrics{}
		data := f.scanner.Bytes()
		err := json.Unmarshal(data, metric)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metric: %w", err)
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
