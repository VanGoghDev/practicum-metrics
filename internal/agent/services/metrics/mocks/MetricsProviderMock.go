package mocks

type MetricsProviderMock struct{}

func (mp MetricsProviderMock) ReadMetrics() *map[string]any {
	return &map[string]any{
		"mockGauge": 12,
	}
}
