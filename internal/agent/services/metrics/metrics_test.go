package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
			logger, _ := zap.NewProduction()
			mp := New(logger)
			metricsMap, _ := mp.ReadMetrics(0)
			assert.NotEmpty(t, metricsMap)
		})
	}
}
