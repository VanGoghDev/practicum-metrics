package metrics_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
)

func TestMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		value      string
	}

	tests := []struct {
		name      string
		params    map[string]string
		gaugesM   map[string]float64
		countersM map[string]int64
		want      want
	}{
		{
			name: "Get gauge",
			params: map[string]string{
				"type": "gauge",
				"name": "test",
			},
			gaugesM: map[string]float64{
				"test": 200,
			},
			want: want{
				statusCode: 200,
				value:      "200",
			},
		},
		{
			name: "Get counter",
			params: map[string]string{
				"type": "counter",
				"name": "testCounter",
			},
			gaugesM: map[string]float64{
				"test": 200,
			},
			countersM: map[string]int64{
				"testCounter": 100,
			},
			want: want{
				statusCode: 200,
				value:      "100",
			},
		},
		// {
		// 	name: "Get empty counter returns 404",
		// 	params: map[string]string{
		// 		"type": "counter",
		// 		"name": "testCounter",
		// 	},
		// 	gaugesM:   map[string]float64{},
		// 	countersM: map[string]int64{},
		// 	want: want{
		// 		statusCode: 404,
		// 		value:      "0",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chirouter.BuildRouter(nil, &mocks.MemStorageMock{
				GaugesM:   tt.gaugesM,
				CountersM: tt.countersM,
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R().SetPathParams(tt.params)
			req.Method = http.MethodGet
			req.URL = fmt.Sprintf("%s/%s", srv.URL, "value/{type}/{name}")
			resp, err := req.Send()
			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.value, resp.String())
		})
	}
}