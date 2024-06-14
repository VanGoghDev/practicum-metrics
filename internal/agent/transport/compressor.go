package transport

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type CompressionTripper struct {
	Proxied http.RoundTripper
}

func (ct *CompressionTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write(body)
	if err != nil {
		return nil, fmt.Errorf("failed to write gzip: %w", err)
	}

	err = zb.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush gzip: %w", err)
	}

	err = zb.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip: %w", err)
	}
	req.Body = io.NopCloser(buf)
	req.ContentLength = int64(buf.Len())
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := ct.Proxied.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to round trip: %w", err)
	}

	return res, nil
}
