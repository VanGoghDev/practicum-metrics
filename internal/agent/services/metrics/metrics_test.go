package metrics

import (
	"context"
	"sync"
	"testing"
	"time"

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
			metricsCh := make(chan Result)
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup

			go mp.ReadMetrics(ctx, metricsCh, time.Second, 1, &wg)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for r := range metricsCh {
					assert.NotEmpty(t, r.Metrics)
				}
			}()
			time.Sleep(200 * time.Millisecond)
			cancel()
			wg.Wait()
		})
	}
}
