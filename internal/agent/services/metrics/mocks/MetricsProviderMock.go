package mocks

type MetricsProviderMock struct{}

func (mp MetricsProviderMock) ReadMetrics() (map[string]any, error) {
	testVal := 12
	return map[string]any{
		"mockGauge": testVal,
	}, nil
}
