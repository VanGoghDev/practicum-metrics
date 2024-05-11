package consumer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrNameIsEmpty      = errors.New("name is empty")
	ErrValueIsIncorrect = errors.New("value is incorrect")
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
		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://%s/update/gauge/%v/%v", s.url, k, v),
			http.NoBody,
		)

		if err != nil {
			return fmt.Errorf("failed to create request for gauge update %w", err)
		}
		resp, err := s.client.Do(request)
		if err != nil {
			return fmt.Errorf("failed to save gauges on server %w", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to close body %w", err)
		}
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
	strValue := strconv.FormatInt(value, 10)
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
	strValue := strconv.FormatFloat(value, 'f', -1, 64)
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
		return fmt.Errorf("failed to save gauge on server %w", err)
	}
	defer dclose(resp.Body)

	return nil
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
