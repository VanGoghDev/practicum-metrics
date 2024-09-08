package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport/signer"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		useSigning bool
	}
	tests := []struct {
		name string
		args args
		want *transport.AgentTripper
	}{
		{
			name: "test signing",
		},
	}
	for _, tt := range tests {
		sgnr := signer.New("secret key")

		t.Run(tt.name, func(t *testing.T) {
			agent := transport.New(tt.args.useSigning, http.DefaultTransport, *sgnr)
			assert.NotNil(t, agent)
		})
	}
}

func TestAgentTripper_RoundTrip(t *testing.T) {
	type args struct {
		useSigning bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "do not use signing",
		},
		{
			name: "use signing",
			args: args{
				useSigning: true,
			},
		},
	}
	for _, tt := range tests {
		sgnr := signer.New("secret key")

		t.Run(tt.name, func(t *testing.T) {
			agent := transport.New(tt.args.useSigning, http.DefaultTransport, *sgnr)

			req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewBufferString("test"))
			res, err := agent.RoundTrip(req)
			defer func() {
				err = res.Body.Close()
			}()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}
