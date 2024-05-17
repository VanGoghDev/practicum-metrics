package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "read metrics",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MetricsProvider{}
			metricsMap, _ := mp.ReadMetrics()
			assert.NotEmpty(t, metricsMap)
		})
	}
}
