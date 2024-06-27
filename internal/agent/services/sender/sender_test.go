package sender_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/http/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/sender"
	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

type args struct {
	name  string
	value float64
}

func tests() []struct {
	name         string
	args         args
	mockHTTPFunc mocks.GetDoFuncType
	err          error
} {
	tests := []struct {
		name         string
		args         args
		mockHTTPFunc mocks.GetDoFuncType
		err          error
	}{
		{
			name: "valid request",
			args: args{
				name:  "test",
				value: 3.1415,
			},
			mockHTTPFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: nil,
		},
		{
			name: "empty name",
			args: args{
				name:  "",
				value: 3.1415,
			},
			mockHTTPFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: nil,
		},
		{
			name: "wrong value",
			args: args{
				name:  "test",
				value: -31.3300412,
			},
			mockHTTPFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: nil,
		},
	}
	return tests
}

func TestSendGauge(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			s := &sender.ServerConsumer{
				Client: &mocks.MockClient{
					DoFunc: tt.mockHTTPFunc,
				},
			}

			metricsV := make([]*models.Metrics, 10)
			metricsV[0] = &models.Metrics{
				ID:    tt.args.name,
				Value: &tt.args.value,
			}

			var wg sync.WaitGroup

			metricsCh := make(chan metrics.Result)
			resultsCh := make(chan sender.Result)

			wg.Add(1)
			go func() {
				metricsCh <- metrics.Result{
					Err:     nil,
					Metrics: metricsV,
				}
				close(metricsCh)
			}()

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			wg.Add(1)
			// Отправим метрики
			// укажем длинный таймаут, чтобы успеть отменить операцию и не идти вторым кругом.
			go s.SendMetrics(ctx, metricsCh, resultsCh, time.Second)

			wg.Add(1)
			// Прочитаем ответ
			go func() {
				for r := range resultsCh {
					assert.Equal(t, tt.err, r.Error)
				}
			}()

			for r := range metricsCh {
				assert.Equal(t, tt.err, r.Err)
				time.Sleep(200 * time.Millisecond)
				cancel()
				return
			}
			wg.Wait()
		})
	}
}
