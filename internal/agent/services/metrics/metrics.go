package metrics

import (
	"reflect"
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
	v := reflect.ValueOf(*m)
	typeV := v.Type()

	metricsMap := make(map[string]any)
	for i := 0; i < v.NumField(); i++ {
		metricsMap[typeV.Field(i).Name] = v.Field(i).Interface()
	}

	return &metricsMap
}
