package app

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/consumer"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
)

var (
	ErrConsumerServiceNil = errors.New("consumer is not initialized")
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

type ServerConsumer interface {
	SendRuntimeGauge(metrics map[string]any) error
	SendCounter(name string, value int64) error
	SendGauge(name string, value float64) error
}

type App struct {
	Consumer        ServerConsumer
	MetricsProvider consumer.MetricsProvider

	reportInterval time.Duration
	pollInterval   time.Duration
}

func New(cfg *config.Config) *App {
	metricsService := metrics.New()
	csmr := consumer.New(metricsService, &http.Client{}, cfg.Address)

	return &App{
		Consumer:        csmr,
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
	if a.Consumer == nil {
		return fmt.Errorf("%s: %w", op, ErrConsumerServiceNil)
	}

	if a.MetricsProvider == nil {
		return fmt.Errorf("%s: %w", op, ErrMetricsProviderNil)
	}

	pollCount := 0
	gauges, err := a.MetricsProvider.ReadMetrics()
	if err != nil {
		return fmt.Errorf("failed to read metrics %s: %w", op, err)
	}

	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			pollCount++
			gauges, err = a.MetricsProvider.ReadMetrics()
			if err != nil {
				return fmt.Errorf("failed to read metrics %w", err)
			}
		case <-reportTicker.C:
			err := a.Consumer.SendRuntimeGauge(gauges)
			if err != nil {
				return fmt.Errorf("failed to send gauges %w", err)
			}

			err = a.Consumer.SendCounter("PollCount", int64(pollCount))
			if err != nil {
				return fmt.Errorf("failed to send count %w", err)
			}
			pollCount = 0

			err = a.Consumer.SendGauge("RandomValue", rand.Float64())
			if err != nil {
				return fmt.Errorf("failed to send random value %w", err)
			}
		}
	}
}
