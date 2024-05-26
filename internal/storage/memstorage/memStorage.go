package memstorage

import (
	"errors"
	"fmt"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"go.uber.org/zap"
)

var (
	ErrGaugesTableNil   = errors.New("gauges table is not initialized")
	ErrCountersTableNil = errors.New("counter table is not initialized")
	ErrNotFound         = errors.New("gauge not found")
)

type MemStorage struct {
	zlog      *zap.Logger
	GaugesM   map[string]float64
	CountersM map[string]int64
}

func New(zlog *zap.Logger) (MemStorage, error) {
	s := MemStorage{
		zlog:      zlog,
		GaugesM:   make(map[string]float64),
		CountersM: make(map[string]int64),
	}

	return s, nil
}

func (s *MemStorage) SaveGauge(name string, value float64) (err error) {
	if s == nil || s.GaugesM == nil {
		return ErrGaugesTableNil
	}

	s.GaugesM[name] = value
	return nil
}

func (s *MemStorage) SaveCount(name string, value int64) (err error) {
	if s == nil || s.CountersM == nil {
		return ErrCountersTableNil
	}

	s.CountersM[name] += value
	return nil
}

func (s *MemStorage) Gauges() (gauges []models.Gauge, err error) {
	if s == nil || s.GaugesM == nil {
		return nil, ErrGaugesTableNil
	}
	gauges = make([]models.Gauge, 0, len(s.GaugesM))
	for k, v := range s.GaugesM {
		gauges = append(gauges, models.Gauge{
			Name:  k,
			Value: v,
		})
	}
	return gauges, err
}

func (s *MemStorage) Counters() (counters []models.Counter, err error) {
	if s == nil || s.CountersM == nil {
		return nil, ErrCountersTableNil
	}

	counters = make([]models.Counter, 0, len(s.CountersM))
	for k, v := range s.CountersM {
		counters = append(counters, models.Counter{
			Name:  k,
			Value: v,
		})
	}
	return counters, err
}

func (s *MemStorage) Gauge(name string) (gauge models.Gauge, err error) {
	if s == nil || s.GaugesM == nil {
		return models.Gauge{}, ErrCountersTableNil
	}
	s.zlog.Info("getting gauge")

	if name == "" {
		return models.Gauge{}, ErrNotFound
	}
	if v, ok := s.GaugesM[name]; !ok {
		s.zlog.Info("gauge does not exist")
		return models.Gauge{}, ErrNotFound
	} else {
		s.zlog.Info("gauge does not exist")

		gauge = models.Gauge{
			Name:  name,
			Value: v,
		}
		return gauge, nil
	}
}

func (s *MemStorage) Counter(name string) (counter models.Counter, err error) {
	if s == nil || s.CountersM == nil {
		return models.Counter{}, ErrCountersTableNil
	}
	s.zlog.Info("getting counter")

	if name == "" {
		s.zlog.Info("name is empty")

		return models.Counter{}, ErrNotFound
	}
	s.zlog.Info("11111111")

	if v, ok := s.CountersM[name]; !ok {
		s.zlog.Info("counter not found")

		return models.Counter{}, ErrNotFound
	} else {
		s.zlog.Info("counter found")

		counter = models.Counter{
			Name:  name,
			Value: v,
		}
		return counter, nil
	}
}

func (s *MemStorage) GetMetrics() ([]*models.Metrics, error) {
	metrics := make([]*models.Metrics, 0)
	for k, v := range s.CountersM {
		metrics = append(metrics, &models.Metrics{
			ID:    k,
			Delta: &v,
			MType: "counter",
		})
	}

	for k, v := range s.GaugesM {
		metrics = append(metrics, &models.Metrics{
			ID:    k,
			Value: &v,
			MType: "gauge",
		})
	}

	return metrics, nil
}

func (s *MemStorage) SaveMetrics(metrics []*models.Metrics) error {
	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			err := s.SaveGauge(v.ID, *v.Value)
			if err != nil {
				return fmt.Errorf("failed to save gauge: %w", err)
			}
		case "counter":
			err := s.SaveCount(v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("failed to save counter: %w", err)
			}
		}
	}
	return nil
}
