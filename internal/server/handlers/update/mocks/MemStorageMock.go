package mocks

type MemStorageMock struct{}

func (s *MemStorageMock) SaveGauge(name string, value float64) (err error) {
	return nil
}

func (s *MemStorageMock) SaveCount(name string, value int64) (err error) {
	return nil
}
