package app

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
}

func New() *App {
	// metrics service
	metricsService := metrics.New()

	consumer := server.New(metricsService, &http.Client{})

	return &App{
		Consumer:        consumer,
		MetricsProvider: metricsService,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
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
	for {
		pollCount++

		if pollCount%5 == 0 {
			err := a.Consumer.SendRuntimeGauge()
			if err != nil {
				return err
			}

			err = a.Consumer.SendCounter("PollCount", int64(pollCount))
			if err != nil {
				return err
			}

			randomValue := randFloats(1.10, 101.98, 5)
			err = a.Consumer.SendGauge("RandomValue", randomValue)
			if err != nil {
				return err
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func randFloats(min, max float64, n int) float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res[0]
}
