package compressor

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// CompressWriter implements http.ResponseWriter.
type CompressWriter struct {
	http.ResponseWriter
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) CompressWriter {
	return CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (cw CompressWriter) Write(data []byte) (int, error) {
	code, err := cw.zw.Write(data)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to %w", err)
	}
	return code, nil
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (cw *CompressWriter) Close() error {
	err := cw.zw.Close()
	if err != nil {
		return fmt.Errorf("failed to close gzip.Writer: %w", err)
	}
	return nil
}

type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to init CompressReader: %w", err)
	}
	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read implements io.ReadCloser.
func (c *CompressReader) Read(p []byte) (n int, err error) {
	code, err := c.zr.Read(p)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to read gzip.Reader: %w", err)
	}
	return code, nil
}

func (c *CompressReader) Close() error {
	err := c.r.Close()
	if err != nil {
		return fmt.Errorf("failed to close io.ReadCloser: %w", err)
	}
	return nil
}

func (c *CompressReader) CloseGzRReader() error {
	err := c.zr.Close()
	if err != nil {
		return fmt.Errorf("failed to close gzip.Reader: %w", err)
	}
	return nil
}

func New(zlog *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			// Если клиент поддерживает обработку сжатых ответов, то переопределим responseWriter.
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				cw := NewCompressWriter(w)
				cw.w.Header().Set("Content-Encoding", "gzip")

				ow = cw
				defer func() {
					err := cw.Close()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}()
			}

			// Если данные пришли в сжатом формате, то заменим body после декомпрессии.
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				zlog.Info("reading compressed body")

				cr, err := NewCompressReader(r.Body)
				if err != nil {
					zlog.Warnf("Unable to create CompressReader: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer func() {
					err = cr.Close()
					if err != nil {
						zlog.Warnf("failed to close compress reader: %v", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}()
			}
			next.ServeHTTP(ow, r)
		}

		return http.HandlerFunc(fn)
	}
}
