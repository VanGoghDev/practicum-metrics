package update_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			log, _ := logger.New("Info")
			s, _ := memstorage.New(log)
			r := chirouter.BuildRouter(s, log, &config.Config{})
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R().SetPathParams(tt.params)
			req.Method = http.MethodPost
			req.URL = fmt.Sprintf("%s/%s", srv.URL, tt.request)
			resp, err := req.Send()

			assert.Empty(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}

// ...

func TestGzipCompression(t *testing.T) {
	log, _ := logger.New("Info")
	s, _ := memstorage.New(log)
	s.GaugesM = map[string]float64{
		"Alloc": 2.0,
	}
	r := chirouter.BuildRouter(s, log, &config.Config{})
	srv := httptest.NewServer(r)
	defer srv.Close()

	requestBody := `{
		"id": "Alloc",
		"type": "gauge",
		"value": 2
    }`

	// ожидаемое содержимое тела ответа при успешном запросе
	successBody := `{
		"id": "Alloc",
		"type": "gauge",
		"value": 2
    }`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		resp, err := resty.New().R().
			SetDoNotParseResponse(true).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "").
			SetBody(buf).
			Post(srv.URL + "/update")

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		b, err := io.ReadAll(resp.RawBody())
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
		require.Equal(t, resp.Header().Get("Accept-Encoding"), "")
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/value", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			err = resp.Body.Close()
			if err != nil {
				require.Empty(t, err)
			}
		}()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)
		require.NotEmpty(t, b)
		require.JSONEq(t, successBody, string(b))
	})
}
