package consumer

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MetricsProvider interface {
	ReadMetrics(pollCount int64) ([]*models.Metrics, error)
}

type ServerConsumer struct {
	metricsProvider MetricsProvider
	client          HTTPClient
	url             string
}

func New(metricsProvider MetricsProvider, client HTTPClient, url string) *ServerConsumer {
	return &ServerConsumer{
		metricsProvider: metricsProvider,
		client:          client,
		url:             url,
	}
}

func (s *ServerConsumer) SendMetrics(metrics []*models.Metrics) error {
	for _, m := range metrics {
		mJ, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to serialize gauge: %w", err)
		}

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err = zb.Write(mJ)
		if err != nil {
			return fmt.Errorf("failed to write gzip: %w", err)
		}
		err = zb.Close()
		if err != nil {
			return fmt.Errorf("failed to close gzip: %w", err)
		}

		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://%s/update/", s.url),
			buf)
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")
		if err != nil {
			return fmt.Errorf("failed to create request for gauge update %w", err)
		}

		err = s.sendRequest(request)
		if err != nil {
			return fmt.Errorf("failed to send request for metric %s: %w", m.ID, err)
		}
	}
	return nil
}

func (s *ServerConsumer) sendRequest(request *http.Request) error {
	resp, err := s.client.Do(request)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return fmt.Errorf("failed to save gauge on server %w", err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("server returned unexpected status code: %d", resp.StatusCode)
	}
	defer func() {
		err = dclose(resp.Body)
	}()

	return err
}

func dclose(c io.Closer) error {
	if err := c.Close(); err != nil {
		return fmt.Errorf("failed to close %w", err)
	}
	return nil
}
