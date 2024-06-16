package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/sender"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport"
	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"go.uber.org/zap"
)

var (
	ErrConsumerServiceNil = errors.New("consumer is not initialized")
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

type Sender interface {
	SendMetrics(metrics []*models.Metrics) error
}

type App struct {
	Log             *zap.Logger
	Sender          Sender
	MetricsProvider sender.MetricsProvider

	reportInterval time.Duration
	pollInterval   time.Duration
}

func New(log *zap.Logger, cfg *config.Config) *App {
	metricsService := metrics.New(log)
	aTripper := transport.New(cfg, http.DefaultTransport)
	sndr := sender.New(
		log,
		metricsService,
		&http.Client{
			Transport: aTripper,
		},
		cfg.Address)

	return &App{
		Log:             log,
		Sender:          sndr,
		MetricsProvider: metricsService,
		reportInterval:  cfg.ReportInterval,
		pollInterval:    cfg.PollInterval,
	}
}

func (a *App) RunApp() error {
	if err := a.Run(); err != nil {
		return fmt.Errorf("failed to run app %w", err)
	}
	return nil
}

func (a *App) Run() error {
	const op = "app.Run"
	if a.Sender == nil {
		return fmt.Errorf("%s: %w", op, ErrConsumerServiceNil)
	}

	if a.MetricsProvider == nil {
		return fmt.Errorf("%s: %w", op, ErrMetricsProviderNil)
	}

	pollCount := 0
	metricsV, err := a.MetricsProvider.ReadMetrics(int64(pollCount))
	if err != nil {
		return fmt.Errorf("failed to read metrics %s: %w", op, err)
	}

	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()
	retriesCount := 0
	maxRetriesCount := 3
	f := 2
	for {
		select {
		case <-pollTicker.C:
			pollCount++
			metricsV, err = a.MetricsProvider.ReadMetrics(int64(pollCount))
			if err != nil {
				a.Log.Warn(fmt.Sprintf("failed to read metrics %s", err))
			}
		case <-reportTicker.C:
			err = a.Sender.SendMetrics(metricsV)
			if err != nil {
				if retriesCount >= maxRetriesCount {
					pollTicker.Stop()
					reportTicker.Stop()
					return fmt.Errorf("tried to send metric %d times, error is: %w", retriesCount, err)
				}
				retriesCount++
				newInterval := (retriesCount * f) - 1
				reportTicker.Stop()
				reportTicker = time.NewTicker(time.Duration(newInterval) * time.Second)
				a.Log.Warn(fmt.Sprintf("failed to send metrics %s", err))
			}
			pollCount = 0
		}
	}
}
