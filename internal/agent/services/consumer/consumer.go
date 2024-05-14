package consumer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/util/converter"
)

var (
	ErrNameIsEmpty        = errors.New("name is empty")
	ErrValueIsIncorrect   = errors.New("value is incorrect")
	ErrServiceUnavailable = errors.New("service is unavailable")
	ErrServiceNotFound    = errors.New("service not found")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MetricsProvider interface {
	ReadMetrics() (map[string]any, error)
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

func (s *ServerConsumer) SendRuntimeGauge(metrics map[string]any) error {
	for k, v := range metrics {
		sV, err := converter.Str(v)
		if err != nil {
			if errors.Is(err, converter.ErrUnsupportedType) {
				continue
			}
			return fmt.Errorf("failed to parse gauge value %w", err)
		}

		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://%s/update/gauge/%v/%v", s.url, k, sV),
			http.NoBody,
		)

		if err != nil {
			return fmt.Errorf("failed to create request for gauge update %w", err)
		}
		return s.sendRequest(request)
	}

	return nil
}

func (s *ServerConsumer) SendCounter(name string, value int64) error {
	if name == "" {
		return ErrNameIsEmpty
	}
	if value < 0 {
		return ErrValueIsIncorrect
	}

	strValue, err := converter.Str(value)
	if err != nil {
		return fmt.Errorf("failed to parse counter value %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/update/counter/%v/%v", s.url, name, strValue),
		http.NoBody)

	if err != nil {
		return fmt.Errorf("failed to create request for counter update %w", err)
	}
	return s.sendRequest(request)
}

func (s *ServerConsumer) SendGauge(name string, value float64) error {
	if name == "" {
		return ErrNameIsEmpty
	}
	if value < 0 {
		return ErrValueIsIncorrect
	}

	strValue, err := converter.Str(value)
	if err != nil {
		return fmt.Errorf("failed to parse gauge value %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/update/gauge/%v/%v", s.url, name, strValue),
		http.NoBody)

	if err != nil {
		return fmt.Errorf("failed to create request for gauge update %w", err)
	}
	return s.sendRequest(request)
}

func (s *ServerConsumer) sendRequest(request *http.Request) error {
	resp, err := s.client.Do(request)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return fmt.Errorf("failed to save gauge on server %w", err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return ErrServiceUnavailable
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrServiceNotFound
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
