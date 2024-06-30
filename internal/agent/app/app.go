package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/services/sender"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport"
	"go.uber.org/zap"
)

var (
	ErrConsumerServiceNil = errors.New("consumer is not initialized")
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

type MetricsProvider interface {
	ReadMetrics(
		ctx context.Context,
		metricsCh chan<- metrics.Result,
		pollInterval time.Duration,
		pollCount int64,
		wg *sync.WaitGroup,
	)
}

type Sender interface {
	SendMetrics(
		ctx context.Context,
		metricsCh <-chan metrics.Result,
		resultCh chan<- sender.Result,
		reportInteval time.Duration,
		wg *sync.WaitGroup,
	)
}

type App struct {
	Log             *zap.Logger
	Sender          Sender
	MetricsProvider MetricsProvider

	rateLimit      int64
	reportInterval time.Duration
	pollInterval   time.Duration
}

func New(log *zap.Logger, cfg *config.Config) *App {
	metricsService := metrics.New(log)
	aTripper := transport.New(cfg, http.DefaultTransport)
	sndr := sender.New(
		log,
		&http.Client{
			Transport: aTripper,
		},
		cfg.Address)

	return &App{
		Log:             log,
		Sender:          sndr,
		MetricsProvider: metricsService,
		rateLimit:       cfg.RateLimit,
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

	var wg sync.WaitGroup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pollCount := 0
	retriesCount := 0
	maxRetriesCount := 3
	_ = maxRetriesCount

	metricsCh := make(chan metrics.Result, a.rateLimit)

	// сюда записываю результат работы "отправителя".
	resultCh := make(chan sender.Result)

	// считаем метрики (горутина спит по установленному таймауту).
	go a.MetricsProvider.ReadMetrics(ctx, metricsCh, a.pollInterval, int64(pollCount), &wg)

	// создадим воркеров и каждый будет отправлять запрос на сервер.
	for w := 1; w <= int(a.rateLimit); w++ {
		go a.Sender.SendMetrics(ctx, metricsCh, resultCh, a.reportInterval, &wg)
	}

	for r := range resultCh {
		if r.Error != nil {
			retriesCount++
			a.Log.Info("result chanel contains errors")
			if maxRetriesCount <= retriesCount {
				cancel()
				return fmt.Errorf("tried to send metric %d times, error is: %w", retriesCount, r.Error)
			}
		}
	}
	wg.Wait()
	return nil
}
