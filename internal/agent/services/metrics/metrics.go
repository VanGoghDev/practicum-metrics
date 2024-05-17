package metrics

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type MetricsProvider struct {
}

func New() *MetricsProvider {
	return &MetricsProvider{}
}

func (mp *MetricsProvider) ReadMetrics() (map[string]any, error) {
	m := new(runtime.MemStats)

	runtime.ReadMemStats(m)

	jM, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json %w", err)
	}

	metricsMap := make(map[string]any)
	err = json.Unmarshal(jM, &metricsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json %w", err)
	}
	return metricsMap, nil
}
