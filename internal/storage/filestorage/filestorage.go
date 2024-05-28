package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/storage"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"go.uber.org/zap"
)

type FileStorage struct {
	memstorage.MemStorage
	zlog    *zap.Logger
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func New(zlog *zap.Logger, cfg *config.Config) (*FileStorage, error) {
	var perm fs.FileMode = 0o666
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to  open a file: %w", err)
	}

	memsrtg, err := memstorage.New(zlog)
	if err != nil {
		return nil, fmt.Errorf("failed to init memory storage: %w", err)
	}

	f := &FileStorage{
		zlog:       zlog,
		file:       file,
		MemStorage: *memsrtg,
		writer:     bufio.NewWriter(file),
		scanner:    bufio.NewScanner(file),
	}
	f.MemStorage.GaugesM = make(map[string]float64)
	f.MemStorage.CountersM = make(map[string]int64)

	if cfg.Restore && f.file != nil {
		err := f.Restore()
		if err != nil {
			return nil, fmt.Errorf("failed to restore file storage: %w", err)
		}
	}
	return f, nil
}

func (f *FileStorage) SaveGauge(name string, value float64) (err error) {
	if f.MemStorage.GaugesM == nil {
		return storage.ErrGaugesTableNil
	}

	f.MemStorage.GaugesM[name] = value

	gauge := &models.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
	data, err := json.Marshal(gauge)
	if err != nil {
		return fmt.Errorf("failed to marshal gauge %s: %w", gauge.ID, err)
	}
	// добавим символ переноса строки
	data = append(data, '\n')

	return f.SaveToFile(data)
}

func (f *FileStorage) SaveCount(name string, value int64) (err error) {
	if f.MemStorage.CountersM == nil {
		return storage.ErrCountersTableNil
	}

	f.MemStorage.CountersM[name] += value

	counter := &models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	}
	data, err := json.Marshal(counter)
	if err != nil {
		return fmt.Errorf("failed to marshal counter %s: %w", counter.ID, err)
	}
	// добавим символ переноса строки
	data = append(data, '\n')

	return f.SaveToFile(data)
}

func (f *FileStorage) SaveToFile(data []byte) error {
	_, err := f.writer.Write(data)

	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = f.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}

	return nil
}

func (f *FileStorage) Restore() error {
	f.zlog.Info("restoring metrics from file!")
	metrics := make([]*models.Metrics, 0)
	for f.scanner.Scan() {
		metric := models.Metrics{}
		data := f.scanner.Bytes()
		if len(data) > 0 {
			f.zlog.Info("data:", zap.String("data:", string(data)))
			err := json.Unmarshal(data, &metric)
			if err != nil {
				fmt.Println(f.scanner.Text())
				return fmt.Errorf("failed to unmarshal metric: %w", err)
			}
			metrics = append(metrics, &metric)
		}
	}
	if err := f.scanner.Err(); err != nil {
		f.zlog.Sugar().Warnf("failed to scan file: %v", err)
	}
	f.zlog.Info("got metrics from file!")

	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			err := f.SaveGauge(v.ID, *v.Value)
			if err != nil {
				return fmt.Errorf("failed to restore gauge %s: %w", v.ID, err)
			}
		case "counter":
			err := f.SaveCount(v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("failed to restore counter %s: %w", v.ID, err)
			}
		}
	}

	return nil
}

func (f *FileStorage) Close() error {
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("Filestorage.Close: %w", err)
	}
	return nil
}
