package consumer

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/http/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
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
			err: ErrNameIsEmpty,
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
			err: ErrValueIsIncorrect,
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
			err := s.SendGauge(tt.args.name, tt.args.value)
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
			err := s.SendCounter(tt.args.name, int64(tt.args.value))
			assert.Equal(t, tt.err, err)
		})
	}
}
