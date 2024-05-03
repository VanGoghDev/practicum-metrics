package mocks

import (
	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

type MemStorageMock struct {
	GaugesM   map[string]float64
	CountersM map[string]int64
}

func (s *MemStorageMock) SaveGauge(name string, value float64) (err error) {
	return nil
}

func (s *MemStorageMock) SaveCount(name string, value int64) (err error) {
	return nil
}

func (s *MemStorageMock) Gauges() (gauges []models.Gauge, err error) {
	return nil, nil
}

func (s *MemStorageMock) Counters() (counters []models.Counter, err error) {
	return nil, nil
}

func (s *MemStorageMock) Gauge(name string) (gauge models.Gauge, err error) {
	return models.Gauge{
		Name:  name,
		Value: s.GaugesM[name],
	}, nil
}

func (s *MemStorageMock) Counter(name string) (counter models.Counter, err error) {
	return models.Counter{
		Name:  name,
		Value: s.CountersM[name],
	}, nil
}
