package consumer

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/http/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
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
			metricsProviderMock := &metrics.MetricsProvider{}
			s := &ServerConsumer{
				metricsProvider: metricsProviderMock,
				client: &mocks.MockClient{
					DoFunc: tt.mockHTTPFunc,
				},
			}
			metricsV := make([]*models.Metrics, 10)
			metricsV[0] = &models.Metrics{
				ID:    tt.args.name,
				Value: &tt.args.value,
			}
			err := s.SendMetrics(metricsV)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSendCounter(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			metricsProviderMock := &metrics.MetricsProvider{}
			s := &ServerConsumer{
				metricsProvider: metricsProviderMock,
				client: &mocks.MockClient{
					DoFunc: tt.mockHTTPFunc,
				},
			}
			metricsV := make([]*models.Metrics, 10)
			val := int64(tt.args.value)
			metricsV[0] = &models.Metrics{
				ID:    tt.args.name,
				Delta: &val,
			}
			err := s.SendMetrics(metricsV)
			assert.Equal(t, tt.err, err)
		})
	}
}
