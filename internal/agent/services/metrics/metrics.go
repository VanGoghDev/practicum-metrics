package metrics

import (
	"encoding/json"
	"runtime"
)

type MetricsProvider struct {
}

func New() *MetricsProvider {
	return &MetricsProvider{}
}

func (mp *MetricsProvider) ReadMetrics() *map[string]any {
	m := new(runtime.MemStats)

	runtime.ReadMemStats(m)

	jM, _ := json.Marshal(m)

	metricsMap := make(map[string]any)
	json.Unmarshal(jM, &metricsMap)

	return &metricsMap
}
