package compressor

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

// CompressBody вызывается, чтобы сжать body запроса.
// Возвращает *http.Request, в котором body сжато
// и проставлены заголовки Content-Encoding: gzip, Accept-Encoding: gzip.
func (ct *CompressionTripper) CompressBody(req *http.Request) (*http.Request, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// ??
	err = req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close request body: %w", err)
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
	return req, nil
}

func (ct *CompressionTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req, err := ct.CompressBody(req)
	if err != nil {
		return nil, fmt.Errorf("failed to compress body: %w", err)
	}

	res, err := ct.Proxied.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to round trip from compression tripper: %w", err)
	}

	return res, nil
}
