package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		gauge   models.Gauge
		request string
		params  map[string]string
		want    want
	}{
		{
			name:    "Valid request",
			request: "update/{type}/{name}/{value}",
			params: map[string]string{
				"type":  "gauge",
				"name":  "test",
				"value": "1",
			},
			gauge: models.Gauge{
				Name:  "test",
				Value: 1,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name:    "Invalid metric name",
			request: "update/{type}/{name}/{value}",
			params: map[string]string{
				"type":  "gauge",
				"name":  "",
				"value": "1",
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "Invalid metric type",
			request: "update/{type}/{name}/{value}",
			params: map[string]string{
				"type":  "guage",
				"name":  "test",
				"value": "1",
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "Invalid metric value",
			request: "update/{type}/{name}/{value}",
			params: map[string]string{
				"type":  "gauge",
				"name":  "test",
				"value": "ss",
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "Invalid url path",
			request: "update/{type}",
			params: map[string]string{
				"type":  "gauge",
				"name":  "",
				"value": "1",
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.HandleFunc(`/update/{type}/{name}/{value}`, UpdateHandler(&memstorage.MemStorage{
				GaugesM: map[string]float64{
					tt.gauge.Name: tt.gauge.Value,
				},
			}))
			srv := httptest.NewServer(r)
			defer srv.Close()
			req := resty.New().R().SetPathParams(tt.params)
			req.Method = http.MethodPost
			req.URL = fmt.Sprintf("%s/%s", srv.URL, "update/{type}/{name}/{value}")
			resp, err := req.Send()

			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}
