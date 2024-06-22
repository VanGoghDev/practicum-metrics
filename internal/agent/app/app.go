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
	ReadMetricsCh(metricsCh chan metrics.Result, pollInterval time.Duration, pollCount int64, rateLimit int64)
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
	a.MetricsProvider.ReadMetricsCh(metricsCh, a.pollInterval, int64(pollCount), a.rateLimit)

	// создадим воркеров и каждый будет отправлять запрос на сервер.
	for w := 1; w <= int(a.rateLimit); w++ {
		wg.Add(1)
		go a.Sender.SendMetricsCh(metricsCh, resultCh, a.reportInterval)
	}

	for r := range resultCh {
		if r.Error != nil {
			retriesCount++
			a.Log.Info("result chanel contains errors")
			// if maxRetriesCount <= retriesCount {
			// 	close(metricsCh)
			// 	close(resultCh)
			// 	return fmt.Errorf("tried to send metric %d times, error is: %w", retriesCount, r.Error)
			// }
		}
	}
	wg.Wait()
	return nil
	// pollTicker := time.NewTicker(a.pollInterval)
	// defer pollTicker.Stop()

	// reportTicker := time.NewTicker(a.reportInterval)
	// defer reportTicker.Stop()
	// retriesCount := 0
	// maxRetriesCount := 3
	// f := 2
	// for {
	// 	select {
	// 	case <-pollTicker.C:
	// 		pollCount++
	// 		a.MetricsProvider.ReadMetricsCh(metricsCh, int64(pollCount), a.rateLimit)
	// 	case <-reportTicker.C:
	// 		// создаем воркеров
	// 		for w := 1; w <= int(a.rateLimit); w++ {
	// 			go a.Sender.SendMetricsCh(metricsCh, resultCh)
	// 			go func() {
	// 				for r := range resultCh {
	// 					if r.Error != nil {
	// 						if retriesCount >= maxRetriesCount {
	// 							pollTicker.Stop()
	// 							reportTicker.Stop()
	// 							close(resultCh)
	// 							return
	// 							//return fmt.Errorf("tried to send metric %d times, worker: %d, error is: %w", retriesCount, w, r.Error)
	// 						}
	// 						retriesCount++
	// 						newInterval := (retriesCount * f) - 1
	// 						reportTicker.Stop()
	// 						reportTicker = time.NewTicker(time.Duration(newInterval) * time.Second)
	// 						a.Log.Warn(fmt.Sprintf("failed to send metrics, worker: %d, error is: %v", w, r.Error))
	// 					}
	// 				}
	// 			}()
	// 		}
	// 		pollCount = 0
	// 	case <-resultCh:
	// 		close(metricsCh)
	// 		a.Log.Info("Канал закрылся")
	// 		return nil
	// 	}
	// }
	// metricsV, err := a.MetricsProvider.ReadMetrics(int64(pollCount))
	// if err != nil {
	// 	return fmt.Errorf("failed to read metrics %s: %w", op, err)
	// }

	// pollTicker := time.NewTicker(a.pollInterval)
	// defer pollTicker.Stop()

	// reportTicker := time.NewTicker(a.reportInterval)
	// defer reportTicker.Stop()
	// retriesCount := 0
	// maxRetriesCount := 3
	// f := 2
	// for {
	// 	select {
	// 	case <-pollTicker.C:
	// 		pollCount++
	// 		metricsV, err = a.MetricsProvider.ReadMetrics(int64(pollCount))
	// 		if err != nil {
	// 			a.Log.Warn(fmt.Sprintf("failed to read metrics %s", err))
	// 		}
	// 	case <-reportTicker.C:
	// 		err = a.Sender.SendMetrics(metricsV)
	// 		if err != nil {
	// 			if retriesCount >= maxRetriesCount {
	// 				pollTicker.Stop()
	// 				reportTicker.Stop()
	// 				return fmt.Errorf("tried to send metric %d times, error is: %w", retriesCount, err)
	// 			}
	// 			retriesCount++
	// 			newInterval := (retriesCount * f) - 1
	// 			reportTicker.Stop()
	// 			reportTicker = time.NewTicker(time.Duration(newInterval) * time.Second)
	// 			a.Log.Warn(fmt.Sprintf("failed to send metrics %s", err))
	// 		}
	// 		pollCount = 0
	// 	}
	// }
}
