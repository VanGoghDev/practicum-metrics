package app

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	server "github.com/VanGoghDev/practicum-metrics/internal/agent/services/consumer/http"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
)

var (
	ErrConsumerServiceNil = errors.New("consumer is not initialized")
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

type App struct {
	Consumer        *server.ServerConsumer
	MetricsProvider *metrics.MetricsProvider

	reportInterval time.Duration
	pollInterval   time.Duration
}

func New(cfg *config.Config) *App {
	metricsService := metrics.New()
	consumer := server.New(metricsService, &http.Client{}, cfg.Address)

	return &App{
		Consumer:        consumer,
		MetricsProvider: metricsService,
		reportInterval:  cfg.ReportInterval,
		pollInterval:    cfg.PollInterval,
	}
}

func (a *App) RunApp() error {
	if err := a.Run(); err != nil {
		return err
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
	metrics := a.MetricsProvider.ReadMetrics()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	sent := make(chan string)
	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-sent:
			pollCount = 0
		case <-pollTicker.C:
			pollCount++
			metrics = a.MetricsProvider.ReadMetrics()
		case <-reportTicker.C:
			err := a.Consumer.SendRuntimeGauge(metrics)
			if err != nil {
				return err
			}

			err = a.Consumer.SendCounter("PollCount", int64(pollCount))
			if err != nil {
				return err
			}

			err = a.Consumer.SendGauge("RandomValue", rand.Float64())
			if err != nil {
				return err
			}
		}
	}
}
