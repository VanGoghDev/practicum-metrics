package server

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	httpMock "github.com/VanGoghDev/practicum-metrics/internal/agent/http/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSendGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name         string
		args         args
		mockHTTPFunc httpMock.GetDoFuncType
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
					StatusCode: 200,
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
					StatusCode: 200,
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
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: ErrValueIsIncorrect,
		},
	}
	for _, tt := range tests {
		httpMock.GetDoFunc = tt.mockHTTPFunc
		t.Run(tt.name, func(t *testing.T) {
			metricsProviderMock := mocks.MetricsProviderMock{}
			s := &ServerConsumer{
				metricsProvider: metricsProviderMock,
				client:          &httpMock.MockClient{},
			}
			err := s.SendGauge(tt.args.name, tt.args.value)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSendCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name         string
		args         args
		mockHttpFunc httpMock.GetDoFuncType
		err          error
	}{
		{
			name: "valid request",
			args: args{
				name:  "test",
				value: 3,
			},
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: nil,
		},
		{
			name: "empty name",
			args: args{
				name:  "",
				value: 3,
			},
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: ErrNameIsEmpty,
		},
		{
			name: "valid request",
			args: args{
				name:  "test",
				value: -31,
			},
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: ErrValueIsIncorrect,
		},
	}
	for _, tt := range tests {
		httpMock.GetDoFunc = tt.mockHttpFunc
		t.Run(tt.name, func(t *testing.T) {
			metricsProviderMock := mocks.MetricsProviderMock{}
			s := &ServerConsumer{
				metricsProvider: metricsProviderMock,
				client:          &httpMock.MockClient{},
			}
			err := s.SendCounter(tt.args.name, tt.args.value)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSendRuntimeGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name         string
		args         args
		mockHttpFunc httpMock.GetDoFuncType
		err          error
	}{
		{
			name: "valid request",
			args: args{
				name:  "test",
				value: 3.1415,
			},
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
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
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
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
			mockHttpFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		httpMock.GetDoFunc = tt.mockHttpFunc
		t.Run(tt.name, func(t *testing.T) {
			metricsProviderMock := mocks.MetricsProviderMock{}
			s := &ServerConsumer{
				metricsProvider: metricsProviderMock,
				client:          &httpMock.MockClient{},
			}
			err := s.SendRuntimeGauge()
			assert.Equal(t, tt.err, err)
		})
	}
}
