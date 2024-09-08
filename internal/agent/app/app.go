// Приложения по сбору и отправки метрик на сервер.
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
	"github.com/VanGoghDev/practicum-metrics/internal/agent/transport/signer"
	"go.uber.org/zap"
)

var (
	// ErrConsumerServiceNil ошибка, что сервис по отправке не инициализирован.
	ErrConsumerServiceNil = errors.New("consumer is not initialized")

	// ErrMetricsProviderNil ошибка, что сервис по чтению метрик не инициализирован.
	ErrMetricsProviderNil = errors.New("metrics provider is not initialized")
)

// MetricsProvider интерфейс, предоставляющий API к сервису чтения метрик.
type MetricsProvider interface {
	// ReadMetrics записывает метрики в metricsCh, опрос метрик происходит по таймауту,
	// определенному в pollInterval.
	ReadMetrics(
		ctx context.Context,
		metricsCh chan<- metrics.Result,
		pollInterval time.Duration,
		pollCount int64,
		wg *sync.WaitGroup,
	)
}

// Sender контракт, которому должен соответствовать сервис отправки метрик.
type Sender interface {
	// SendMetrics отправляет метрики в resultCh, отправка метрик происходит по таймауту reportInterval.
	SendMetrics(
		ctx context.Context,
		metricsCh <-chan metrics.Result,
		resultCh chan<- sender.Result,
		reportInteval time.Duration,
		wg *sync.WaitGroup,
	)
}

// App ядро приложения. Инкапсулирует в себе сервисы чтения, отправки метрик.
type App struct {
	// Log Логгер
	Log *zap.Logger
	// Sender Сервис отправки метрик.
	Sender Sender

	// MetricsProvider Сервис
	MetricsProvider MetricsProvider

	rateLimit      int64
	reportInterval time.Duration
	pollInterval   time.Duration
}

// New возвращает новый экземпляр приложения.
func New(log *zap.Logger, cfg *config.Config) *App {
	metricsService := metrics.New(log)
	sgnr := signer.New(cfg.Key)
	var useSigning bool
	if cfg.Key != "" {
		useSigning = true
	}
	aTripper := transport.New(useSigning, http.DefaultTransport, *sgnr)
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

// RunApp запускает приложение.
func (a *App) RunApp() error {
	if err := a.Run(); err != nil {
		return fmt.Errorf("failed to run app %w", err)
	}
	return nil
}

// Run запускает приложение.
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
