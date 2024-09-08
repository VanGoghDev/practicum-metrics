package compressor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressionTripper_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		want    *http.Response
		wantErr bool
	}{
		{
			name: "compress body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &CompressionTripper{
				Proxied: http.DefaultTransport,
			}
			req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewBufferString("test"))

			res, err := ct.RoundTrip(req)
			defer func() {
				err = res.Body.Close()
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("CompressionTripper.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
