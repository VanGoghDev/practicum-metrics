package metrics

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestReadMetrics(t *testing.T) {
	tests := []struct {
		name               string
		pollCount          int
		contextFinishAfter int
	}{
		{
			name:               "context finished faster than read metrics",
			pollCount:          20,
			contextFinishAfter: 10,
		},
		{
			name:               "read metrics faster than context finished",
			pollCount:          10,
			contextFinishAfter: 20,
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
			go mp.ReadMetrics(ctx, metricsCh, time.Millisecond*time.Duration(tt.pollCount), 1, &wg)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for r := range metricsCh {
					assert.NotEmpty(t, r.Metrics)
				}
			}()
			time.Sleep(time.Duration(tt.contextFinishAfter) * time.Millisecond)
			cancel()
			wg.Wait()
		})
	}
}

func BenchmarkReadMetrics(b *testing.B) {
	logger, _ := zap.NewProduction()
	mp := New(logger)
	metricsCh := make(chan Result, 10)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var wg sync.WaitGroup
	b.ResetTimer()

	b.Run("read metrics", func(b *testing.B) {
		mp.ReadMetrics(ctx, metricsCh, time.Second, 1, &wg)
	})
}

func Test_populateMetric(t *testing.T) {
	type args struct {
		metric    *models.Metrics
		pollCount int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "populate Alloc",
			args: args{
				metric: &models.Metrics{
					ID: "Alloc",
				},
			},
		},
		{
			name: "populate BuckHashSys",
			args: args{
				metric: &models.Metrics{
					ID: "BuckHashSys",
				},
			},
		},
		{
			name: "populate Frees",
			args: args{
				metric: &models.Metrics{
					ID: "Frees",
				},
			},
		},
		{
			name: "populate GCCPUFraction",
			args: args{
				metric: &models.Metrics{
					ID: "GCCPUFraction",
				},
			},
		},
		{
			name: "populate GCSys",
			args: args{
				metric: &models.Metrics{
					ID: "GCSys",
				},
			},
		},
		{
			name: "populate HeapAlloc",
			args: args{
				metric: &models.Metrics{
					ID: "HeapAlloc",
				},
			},
		},
		{
			name: "populate HeapSys",
			args: args{
				metric: &models.Metrics{
					ID: "HeapSys",
				},
			},
		},
		{
			name: "populate HeapIdle",
			args: args{
				metric: &models.Metrics{
					ID: "HeapIdle",
				},
			},
		},
		{
			name: "populate HeapInuse",
			args: args{
				metric: &models.Metrics{
					ID: "HeapInuse",
				},
			},
		},
		{
			name: "populate HeapObjects",
			args: args{
				metric: &models.Metrics{
					ID: "HeapObjects",
				},
			},
		},
		{
			name: "populate HeapReleased",
			args: args{
				metric: &models.Metrics{
					ID: "HeapReleased",
				},
			},
		},
		{
			name: "populate LastGC",
			args: args{
				metric: &models.Metrics{
					ID: "LastGC",
				},
			},
		},
		{
			name: "populate Lookups",
			args: args{
				metric: &models.Metrics{
					ID: "Lookups",
				},
			},
		},
		{
			name: "populate MCacheInuse",
			args: args{
				metric: &models.Metrics{
					ID: "MCacheInuse",
				},
			},
		},
		{
			name: "populate MCacheSys",
			args: args{
				metric: &models.Metrics{
					ID: "MCacheSys",
				},
			},
		},
		{
			name: "populate MSpanInuse",
			args: args{
				metric: &models.Metrics{
					ID: "MSpanInuse",
				},
			},
		},
		{
			name: "populate MSpanSys",
			args: args{
				metric: &models.Metrics{
					ID: "MSpanSys",
				},
			},
		},

		{
			name: "populate Mallocs",
			args: args{
				metric: &models.Metrics{
					ID: "Mallocs",
				},
			},
		}, {
			name: "populate NextGC",
			args: args{
				metric: &models.Metrics{
					ID: "NextGC",
				},
			},
		},
		{
			name: "populate NumForcedGC",
			args: args{
				metric: &models.Metrics{
					ID: "NumForcedGC",
				},
			},
		},
		{
			name: "populate NumGC",
			args: args{
				metric: &models.Metrics{
					ID: "NumGC",
				},
			},
		},
		{
			name: "populate OtherSys",
			args: args{
				metric: &models.Metrics{
					ID: "OtherSys",
				},
			},
		},
		{
			name: "populate PauseTotalNs",
			args: args{
				metric: &models.Metrics{
					ID: "PauseTotalNs",
				},
			},
		},
		{
			name: "populate StackInuse",
			args: args{
				metric: &models.Metrics{
					ID: "StackInuse",
				},
			},
		},
		{
			name: "populate StackSys",
			args: args{
				metric: &models.Metrics{
					ID: "StackSys",
				},
			},
		},
		{
			name: "populate Sys",
			args: args{
				metric: &models.Metrics{
					ID: "Sys",
				},
			},
		},
		{
			name: "populate TotalAlloc",
			args: args{
				metric: &models.Metrics{
					ID: "TotalAlloc",
				},
			},
		},
		{
			name: "populate TotalMemory",
			args: args{
				metric: &models.Metrics{
					ID: "TotalMemory",
				},
			},
		},
		{
			name: "populate FreeMemory",
			args: args{
				metric: &models.Metrics{
					ID: "FreeMemory",
				},
			},
		},
		{
			name: "populate CPUutilization1",
			args: args{
				metric: &models.Metrics{
					ID: "CPUutilization1",
				},
			},
		},
		{
			name: "populate RandomValue",
			args: args{
				metric: &models.Metrics{
					ID: "RandomValue",
				},
			},
		},
		{
			name: "populate PollCount",
			args: args{
				metric: &models.Metrics{
					ID: "PollCount",
				},
			},
		},
		{
			name: "populate unknown metric returns error",
			args: args{
				metric: &models.Metrics{
					ID: "unknown",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		m := new(runtime.MemStats)
		vm, _ := mem.VirtualMemory()
		t.Run(tt.name, func(t *testing.T) {
			if err := populateMetric(tt.args.metric, m, vm, tt.args.pollCount); (err != nil) != tt.wantErr {
				t.Errorf("populateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
