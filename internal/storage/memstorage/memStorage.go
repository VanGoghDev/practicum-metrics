package memStorage

import (
	"fmt"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func New() (MemStorage, error) {
	s := MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	return s, nil
}

func (s *MemStorage) SaveGauge(name string, value float64) (err error) {
	if s == nil || s.Gauges == nil {
		return fmt.Errorf("gauges table is not initialized")
	}

	s.Gauges[name] = value
	return nil
}

func (s *MemStorage) SaveCount(name string, value int64) (err error) {
	if s == nil || s.Counters == nil {
		return fmt.Errorf("counters table is not initialized")
	}

	if _, ok := s.Counters[name]; ok {
		s.Counters[name] = s.Counters[name] + value
	} else {
		s.Counters[name] = value
	}
	return nil
}
