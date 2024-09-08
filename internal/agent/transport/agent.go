package transport

import (
	"fmt"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport/compressor"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport/signer"
)

// AgentTripper инкапсулирует в себе логику
// транспорта сжатия (CompressionTripper)
// и транспорта подписи (SignTripper).
type AgentTripper struct {
	Proxied http.RoundTripper

	compressor.CompressionTripper
	signer.SignerTripper

	useCompression, useSigning bool
}

// New возвращает новый экземпляр AgentTripper.
func New(useSigning bool, proxy http.RoundTripper, sgnr signer.SignerTripper) *AgentTripper {
	return &AgentTripper{
		SignerTripper:  sgnr,
		useCompression: true,
		useSigning:     useSigning,
		Proxied:        proxy,
	}
}

// RoundTrip сжимает запросы и подписывает их через SignerTripper.
func (a *AgentTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error

	if a.useCompression {
		req, err = a.CompressionTripper.CompressBody(req)
		if err != nil {
			return nil, fmt.Errorf("failed to compress request body: %w", err)
		}
	}

	if a.useSigning {
		req, err = a.SignerTripper.SignBody(req)
		if err != nil {
			return nil, fmt.Errorf("failed to sign request body: %w", err)
		}
	}

	res, err := a.Proxied.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to round trip from agent tripper: %w", err)
	}

	return res, nil
}
