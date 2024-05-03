package server

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNameIsEmpty      = errors.New("name is empty")
	ErrValueIsIncorrect = errors.New("value is incorrect")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MetricsProvider interface {
	ReadMetrics() *map[string]any
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

// sends gauge taken from metrics provider (runtime metrics)
func (s *ServerConsumer) SendRuntimeGauge(metrics *map[string]any) error {
	for k, v := range *metrics {
		request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/gauge/%v/%v", s.url, k, v), nil)
		if err != nil {
			return err
		}
		resp, err := s.client.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		//_, err := http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", k, v), "text/plain", nil)
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
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/counter/%v/%v", s.url, name, value), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *ServerConsumer) SendGauge(name string, value float64) error {
	if name == "" {
		return ErrNameIsEmpty
	}
	if value < 0 {
		return ErrValueIsIncorrect
	}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/gauge/%v/%v", s.url, name, value), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
