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
}

func New(metricsProvider MetricsProvider, client HTTPClient) *ServerConsumer {
	return &ServerConsumer{
		metricsProvider: metricsProvider,
		client:          client,
	}
}

// sends gauge taken from metrics provider (runtime metrics)
func (s *ServerConsumer) SendRuntimeGauge() error {
	metrics := s.metricsProvider.ReadMetrics()

	for k, v := range *metrics {
		request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", k, v), nil)
		if err != nil {
			return err
		}
		resp, err := s.client.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		//_, err := http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", k, v), "text/plain", nil)
		if err != nil {
			return err
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
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/counter/%v/%v", name, value), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerConsumer) SendGauge(name string, value float64) error {
	if name == "" {
		return ErrNameIsEmpty
	}
	if value < 0 {
		return ErrValueIsIncorrect
	}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", name, value), nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
