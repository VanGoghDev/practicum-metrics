package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update/mocks"
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
		request string
		want    want
	}{
		{
			name:    "Valid request",
			request: "gauge/test/1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name:    "Invalid metric name",
			request: "gauge//1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "Invalid metric type",
			request: "guage/test/1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "Invalid metric value",
			request: "gauge/test/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "Invalid url path",
			request: "/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%v", tt.request), nil)
			w := httptest.NewRecorder()

			UpdateHandler(&mocks.MemStorageMock{})(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
