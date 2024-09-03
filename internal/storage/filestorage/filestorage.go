package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"
	"go.uber.org/zap"
)

// Filewriter интерфейс работы с файлом.
type Filewriter interface {
	// SaveMetrics сохраняет метрики в файл.
	SaveMetrics(ctx context.Context, data []byte) error

	// ReadMetrics возвращает метрики, считанные из файла.
	ReadMetrics() ([]*models.Metrics, error)

	// FileIsNil возвращает true если файл nil
	FileIsNil() bool

	// Close закрывает файл.
	Close() error
}

// FileStorage хранилище метрик в файле ОС.
type FileStorage struct {
	memstorage.MemStorage
	zlog *zap.Logger

	filewriter Filewriter
}

// New возвращает новый экземпляр хранилища.
func New(
	ctx context.Context,
	zlog *zap.Logger,
	restore bool,
	fwriter Filewriter,
	memstrg *memstorage.MemStorage,
) (*FileStorage, error) {
	if memstrg == nil {
		return nil, errors.New("mem storage is nil")
	}

	if fwriter.FileIsNil() {
		return nil, errors.New("file is nil")
	}

	f := &FileStorage{
		zlog:       zlog,
		filewriter: fwriter,
		MemStorage: *memstrg,
	}

	if restore {
		err := f.restore(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to restore file storage: %w", err)
		}
	}

	return f, nil
}

// SaveMetrics сохраняет метрики.
func (f *FileStorage) SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error) {
	for _, v := range metrics {
		switch v.MType {
		case handlers.Counter:
			err := f.SaveCount(ctx, v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("failed to save count: %w", err)
			}
		case handlers.Gauge:
			err := f.SaveGauge(ctx, v.ID, *v.Value)
			if err != nil {
				return fmt.Errorf("failed to save gauge: %w", err)
			}
		default:
			f.zlog.Sugar().Warnf("metric \"%s\" has unknown type \"%s\".", v.MType, v.MType)
			continue
		}
	}

	return nil
}

// SaveGauge сохраняет метрики типа Gauge.
func (f *FileStorage) SaveGauge(ctx context.Context, name string, value float64) (err error) {
	if f.MemStorage.GaugesM == nil {
		return serrors.ErrGaugesTableNil
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

	err = f.filewriter.SaveMetrics(ctx, data)
	if err != nil {
		return fmt.Errorf("%w: failed to save metrics to file", err)
	}
	return nil
}

// SaveCount сохраняет метрики типа Count.
func (f *FileStorage) SaveCount(ctx context.Context, name string, value int64) (err error) {
	if f.MemStorage.CountersM == nil {
		return serrors.ErrCountersTableNil
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

	err = f.filewriter.SaveMetrics(ctx, data)
	if err != nil {
		return fmt.Errorf("%w: failed to save metrics to file", err)
	}
	return nil
}

func (f *FileStorage) restore(ctx context.Context) error {
	f.zlog.Debug("restoring metrics from file...")
	metrics, err := f.filewriter.ReadMetrics()
	if err != nil {
		f.zlog.Sugar().Warnf("failed to read metrics from file: %v", err)
	}

	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			err := f.SaveGauge(ctx, v.ID, *v.Value)
			if err != nil {
				return fmt.Errorf("failed to restore gauge %s: %w", v.ID, err)
			}
		case "counter":
			err := f.SaveCount(ctx, v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("failed to restore counter %s: %w", v.ID, err)
			}
		}
	}

	return nil
}

// Ping пинг хранилища.
func (f *FileStorage) Ping(ctx context.Context) error {
	return nil
}

// Close закрывает соединение с хранилищем.
func (f *FileStorage) Close(ctx context.Context) error {
	if err := f.filewriter.Close(); err != nil {
		return fmt.Errorf("Filestorage.Close: %w", err)
	}
	return nil
}
