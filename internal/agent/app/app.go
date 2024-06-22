package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
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

type MetricsProvider interface {
	ReadMetricsCh(metricsCh chan metrics.Result, pollInterval time.Duration, pollCount int64)
	ReadAdditionalMetrics(metricsCh chan metrics.Result, pollInterval time.Duration)
	ReadMetrics(pollCount int64) ([]*models.Metrics, error)
}

type Sender interface {
	SendMetricsCh(metricsCh chan metrics.Result, resultCh chan sender.Result, reportInteval time.Duration)
	SendMetrics(metrics []*models.Metrics) error
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

	pollCount := 0
	retriesCount := 0
	maxRetriesCount := 3
	_ = maxRetriesCount
	// сюда записываю метрики
	// не знаю какой канал использовать правильнее. Буферизированный или нет.
	// Мне показалось, что лучше небуферизированный, ведь операция чтения метрики гораздо быстрее отправки.
	metricsCh := make(chan metrics.Result)

	// сюда записываю результат работы "отправителя".
	resultCh := make(chan sender.Result)

	wg.Add(1)
	// считаем метрики (горутина спит по установленному таймауту).
	a.MetricsProvider.ReadMetricsCh(metricsCh, a.pollInterval, int64(pollCount))
	a.MetricsProvider.ReadAdditionalMetrics(metricsCh, a.pollInterval)

	// создадим воркеров и каждый будет отправлять запрос на сервер.
	for w := 1; w <= int(a.rateLimit); w++ {
		wg.Add(1)
		go a.Sender.SendMetricsCh(metricsCh, resultCh, a.reportInterval)
	}

	for r := range resultCh {
		if r.Error != nil {
			retriesCount++
			a.Log.Info("result chanel contains errors")
			if maxRetriesCount <= retriesCount {
				// Мне кажется я совершенно неверно закрываю каналы. Как я понял, закрывать нужно
				// там же где ты и пишешь. Вот не знаю, надо было как-то через контекст реализовать это?
				// Например если в генерилку передавать контекст и отменять его тут.
				close(metricsCh)
				close(resultCh)
				return fmt.Errorf("tried to send metric %d times, error is: %w", retriesCount, r.Error)
			}
		}
	}
	wg.Wait()
	return nil
}
