package metrics_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	mock_routers "github.com/VanGoghDev/practicum-metrics/internal/server/routers/mocks"
)

func TestMetricsHandler(t *testing.T) {
	type args struct {
		gaugesErr   bool
		countersErr bool
	}
	type want struct {
		statusCode int
		value      string
	}

	tests := []struct {
		name      string
		args      args
		params    map[string]string
		gaugesM   []models.Gauge
		countersM []models.Counter
		want      want
	}{
		{
			name: "storage returns errors when accessing gauges, should return internal server error",
			args: args{
				gaugesErr:   true,
				countersErr: false,
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   []models.Gauge{{Name: "test", Value: 3.14}},
			countersM: []models.Counter{{Name: "test", Value: 3}},
		},
		{
			name: "storage returns errors when accessing counters, should return internal server error",
			args: args{
				gaugesErr:   false,
				countersErr: true,
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   []models.Gauge{{Name: "test", Value: 3.14}},
			countersM: []models.Counter{{Name: "test", Value: 3}},
		},
		{
			name: "storage returns errors when accessing counters, should return internal server error",
			args: args{
				gaugesErr:   false,
				countersErr: false,
			},
			want: want{
				statusCode: http.StatusOK,
			},
			gaugesM:   []models.Gauge{{Name: "test", Value: 3.14}},
			countersM: []models.Counter{{Name: "test", Value: 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_routers.NewMockStorage(ctrl)
			if tt.args.gaugesErr {
				s.EXPECT().Gauges(gomock.Any()).Return(nil, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Gauges(gomock.Any()).Return(tt.gaugesM, nil).AnyTimes()
			}

			if tt.args.countersErr {
				s.EXPECT().Counters(gomock.Any()).Return(nil, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Counters(gomock.Any()).Return(tt.countersM, nil).AnyTimes()
			}

			log, _ := logger.New("Info")
			r := chirouter.BuildRouter(s, log, &config.Config{})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R().SetPathParams(tt.params)
			req.Method = http.MethodGet
			req.URL = fmt.Sprintf("%s/%s", srv.URL, "")
			resp, err := req.Send()
			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.args.countersErr || tt.args.gaugesErr {
				assert.Equal(t, tt.want.value, resp.String())
			}
		})
	}
}

func TestMetricHandlerPost(t *testing.T) {
	type args struct {
		gaugesErr   bool
		countersErr bool
		body        string
	}
	type want struct {
		statusCode int
		value      string
	}

	tests := []struct {
		name      string
		args      args
		params    map[string]string
		gaugesM   models.Gauge
		countersM models.Counter
		want      want
	}{
		{
			name: "storage returns errors when accessing counters, should return internal server error",
			args: args{
				gaugesErr:   false,
				countersErr: true,
				body:        "{\"id\": \"cc2\",\"type\": \"counter\" }",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{Name: "test", Value: 3.14},
			countersM: models.Counter{Name: "test", Value: 3},
		},
		{
			name: "storage returns errors when accessing gauges, should return internal server error",
			args: args{
				gaugesErr:   true,
				countersErr: false,
				body:        "{\"id\": \"cc2\",\"type\": \"gauge\" }",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{Name: "test", Value: 3.14},
			countersM: models.Counter{Name: "test", Value: 3},
		},
		{
			name: "returns ok gauge",
			args: args{
				gaugesErr:   false,
				countersErr: false,
				body:        "{\"id\": \"cc2\",\"type\": \"gauge\" }",
			},
			want: want{
				statusCode: http.StatusOK,
			},
			gaugesM:   models.Gauge{Name: "test", Value: 3.14},
			countersM: models.Counter{Name: "test", Value: 3},
		},
		{
			name: "returns ok counter",
			args: args{
				gaugesErr:   false,
				countersErr: false,
				body:        "{\"id\": \"cc2\",\"type\": \"counter\" }",
			},
			want: want{
				statusCode: http.StatusOK,
			},
			gaugesM:   models.Gauge{Name: "test", Value: 3.14},
			countersM: models.Counter{Name: "test", Value: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_routers.NewMockStorage(ctrl)
			if tt.args.gaugesErr {
				s.EXPECT().Gauge(gomock.Any(), gomock.Any()).Return(tt.gaugesM, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Gauge(gomock.Any(), gomock.Any()).Return(tt.gaugesM, nil).AnyTimes()
			}

			if tt.args.countersErr {
				s.EXPECT().Counter(gomock.Any(), gomock.Any()).Return(tt.countersM, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Counter(gomock.Any(), gomock.Any()).Return(tt.countersM, nil).AnyTimes()
			}

			log, _ := logger.New("Info")
			r := chirouter.BuildRouter(s, log, &config.Config{})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R().SetBody(bytes.NewBufferString(tt.args.body))

			req.Method = http.MethodPost
			req.URL = fmt.Sprintf("%s/%s", srv.URL, "value")
			resp, err := req.Send()
			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.args.countersErr || tt.args.gaugesErr {
				assert.Equal(t, tt.want.value, resp.String())
			}
		})
	}
}

func TestMetricHandler(t *testing.T) {
	type args struct {
		gaugesErr   bool
		countersErr bool
		urltype     string
		urlName     string
	}
	type want struct {
		statusCode int
		value      string
	}

	tests := []struct {
		name      string
		args      args
		gaugesM   models.Gauge
		countersM models.Counter
		want      want
	}{
		{
			name: "Invalid url type param returns bad request",
			args: args{
				gaugesErr: true,
				urltype:   "invalid",
				urlName:   "test",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				value:      "Invalid metric type",
			},
			gaugesM:   models.Gauge{},
			countersM: models.Counter{},
		},
		{
			name: "Storage returns error when get gauge, should return internal error",
			args: args{
				gaugesErr: true,
				urltype:   "gauge",
				urlName:   "test",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{},
			countersM: models.Counter{},
		},
		{
			name: "returns ok when get gauge",
			args: args{
				urltype: "gauge",
				urlName: "test",
			},
			want: want{
				statusCode: http.StatusOK,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{},
			countersM: models.Counter{},
		},
		{
			name: "Storage returns error when get counter, should return internal error",
			args: args{
				countersErr: true,
				urltype:     "counter",
				urlName:     "test",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{},
			countersM: models.Counter{},
		},
		{
			name: "returns ok when get counter",
			args: args{
				urltype: "counter",
				urlName: "test",
			},
			want: want{
				statusCode: http.StatusOK,
				value:      "Internal error",
			},
			gaugesM:   models.Gauge{},
			countersM: models.Counter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, _ := logger.New("Info")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_routers.NewMockStorage(ctrl)
			if tt.args.gaugesErr {
				s.EXPECT().Gauge(gomock.Any(), gomock.Any()).Return(tt.gaugesM, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Gauge(gomock.Any(), gomock.Any()).Return(tt.gaugesM, nil).AnyTimes()
			}

			if tt.args.countersErr {
				s.EXPECT().Counter(gomock.Any(), gomock.Any()).Return(tt.countersM, errors.New("test error")).AnyTimes()
			} else {
				s.EXPECT().Counter(gomock.Any(), gomock.Any()).Return(tt.countersM, nil).AnyTimes()
			}
			r := chirouter.BuildRouter(s, log, &config.Config{})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = fmt.Sprintf("%s/value/%s/%s", srv.URL, tt.args.urltype, tt.args.urlName)
			resp, err := req.Send()
			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.args.countersErr || tt.args.gaugesErr {
				assert.Equal(t, tt.want.value, resp.String())
			}
		})
	}
}
