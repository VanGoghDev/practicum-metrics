package memstorage

import (
	"errors"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

var (
	ErrGaugesTableNil   = errors.New("gauges table is not initialized")
	ErrCountersTableNil = errors.New("counter table is not initialized")
	ErrNotFound         = errors.New("gauge not found")
)

type MemStorage struct {
	GaugesM   map[string]float64
	CountersM map[string]int64
}

func New() (MemStorage, error) {
	s := MemStorage{
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

	if _, ok := s.CountersM[name]; ok {
		s.CountersM[name] = s.CountersM[name] + value
	} else {
		s.CountersM[name] = value
	}
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
	return
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
	return
}

func (s *MemStorage) Gauge(name string) (gauge models.Gauge, err error) {
	if s == nil || s.CountersM == nil {
		return models.Gauge{}, ErrCountersTableNil
	}

	if name == "" {
		return models.Gauge{}, ErrNotFound
	}

	if v, ok := s.GaugesM[name]; !ok {
		return models.Gauge{}, ErrNotFound
	} else {
		gauge = models.Gauge{
			Name:  name,
			Value: v,
		}
		return
	}
}

func (s *MemStorage) Counter(name string) (counter models.Counter, err error) {
	if s == nil || s.CountersM == nil {
		return models.Counter{}, ErrCountersTableNil
	}

	if name == "" {
		return models.Counter{}, ErrNotFound
	}

	if v, ok := s.CountersM[name]; !ok {
		return models.Counter{}, ErrNotFound
	} else {
		counter = models.Counter{
			Name:  name,
			Value: v,
		}
		return
	}
}
