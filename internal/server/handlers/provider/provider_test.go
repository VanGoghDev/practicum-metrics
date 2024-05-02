package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/mocks"
)

func TestMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		params  map[string]string
		gaugesM map[string]float64
		want    want
	}{
		{
			name: "Valid request",
			params: map[string]string{
				"type": "gauge",
				"name": "test",
			},
			gaugesM: map[string]float64{
				"test": 200,
			},
			want: want{
				statusCode: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Route("/value", func(r chi.Router) {
				r.Get("/{type}/{name}", MetricHandler(&mocks.MemStorageMock{
					GaugesM: tt.gaugesM,
				}))
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R().SetPathParams(tt.params)
			req.Method = http.MethodGet
			req.URL = fmt.Sprintf("%s/%s", srv.URL, "value/{type}/{name}")
			resp, err := req.Send()
			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
		})
	}
}
