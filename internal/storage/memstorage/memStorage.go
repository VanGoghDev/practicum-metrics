package memstorage

import (
	"errors"
	"fmt"
)

var (
	ErrGaugesTableNil   = errors.New("gauges table is not initialized")
	ErrCountersTableNil = errors.New("counter table is not initialized")
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
		return ErrGaugesTableNil
	}

	s.Gauges[name] = value
	fmt.Printf("%v %v /n", name, value)
	return nil
}

func (s *MemStorage) SaveCount(name string, value int64) (err error) {
	if s == nil || s.Counters == nil {
		return ErrCountersTableNil
	}

	if _, ok := s.Counters[name]; ok {
		s.Counters[name] = s.Counters[name] + value
	} else {
		s.Counters[name] = value
	}
	return nil
}
