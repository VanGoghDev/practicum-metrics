package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/consumer"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"go.uber.org/zap"
)

var (
	ErrConsumerServiceNil = errors.New("consumer is not initialized")
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

type ServerConsumer interface {
	SendMetrics(metrics []*models.Metrics) error
}

type App struct {
	Log             *zap.Logger
	Consumer        ServerConsumer
	MetricsProvider consumer.MetricsProvider

	reportInterval time.Duration
	pollInterval   time.Duration
}

func New(log *zap.Logger, cfg *config.Config) *App {
	metricsService := metrics.New(log)
	csmr := consumer.New(metricsService, &http.Client{}, cfg.Address)

	return &App{
		Log:             log,
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
	metricsV, err := a.MetricsProvider.ReadMetrics(int64(pollCount))
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
			metricsV, err = a.MetricsProvider.ReadMetrics(int64(pollCount))
			if err != nil {
				a.Log.Warn(fmt.Sprintf("failed to read metrics %s", err))
			}
		case <-reportTicker.C:
			err = a.Consumer.SendMetrics(metricsV)
			if err != nil {
				a.Log.Warn(fmt.Sprintf("failed to send metrics %s", err))
			}
			pollCount = 0
		}
	}
}
