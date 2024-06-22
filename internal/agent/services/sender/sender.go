package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"go.uber.org/zap"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Result struct {
	Error error
}

type ServerConsumer struct {
	zlog   *zap.Logger
	client HTTPClient
	url    string
}

func New(zlog *zap.Logger, client HTTPClient, url string) *ServerConsumer {
	return &ServerConsumer{
		zlog:   zlog,
		client: client,
		url:    url,
	}
}

func (s *ServerConsumer) SendMetricsCh(
	metricsCh chan metrics.Result,
	resultCh chan Result,
	reportInteval time.Duration) {
	// defer close(resultCh) //// как закрыть канал? Может через сигнал?
	for {
		s.zlog.Info("отправляем метрики")
		for m := range metricsCh {
			if m.Err != nil {
				s.zlog.Warn(fmt.Sprintf("failed to read metrics %v", m.Err))
				continue
			}
			mJ, err := json.Marshal(m.Metrics)
			if err != nil {
				s.zlog.Warn(fmt.Sprintf("failed to serialize gauge: %v", err))
				resultCh <- Result{
					Error: err,
				}
				continue
			}

			buf := bytes.NewBuffer(mJ)

			request, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("http://%s/updates/", s.url),
				buf)
			request.Close = true
			if err != nil {
				s.zlog.Warn(fmt.Sprintf("failed to create request for metrics update: %v", err))
				resultCh <- Result{
					Error: err,
				}
				continue
			}

			err = s.sendRequest(request)
			if err != nil {
				s.zlog.Warn(fmt.Sprintf("failed to send request: %v", err))
				resultCh <- Result{
					Error: err,
				}
				continue
			}
		}
		time.Sleep(reportInteval)
	}

}

func (s *ServerConsumer) SendMetrics(mtrcs []*models.Metrics) error {
	mJ, err := json.Marshal(mtrcs)
	if err != nil {
		return fmt.Errorf("failed to serialize gauge: %w", err)
	}

	buf := bytes.NewBuffer(mJ)

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/updates/", s.url),
		buf)
	request.Close = true
	if err != nil {
		return fmt.Errorf("failed to create request for metrics update %w", err)
	}

	err = s.sendRequest(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	return nil
}

func (s *ServerConsumer) sendRequest(request *http.Request) error {
	resp, err := s.client.Do(request)
	if err != nil {
		s.zlog.Sugar().Errorf("failed to send request: %w", err)
		return fmt.Errorf("failed to send metrics to server %w", err)
	}
	defer func() {
		err = dclose(resp.Body)
	}()

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("server returned unexpected status code: %d", resp.StatusCode)
	}

	return err
}

func dclose(c io.Closer) error {
	if err := c.Close(); err != nil {
		return fmt.Errorf("failed to close %w", err)
	}
	return nil
}
