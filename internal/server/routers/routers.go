package routers

import (
	"context"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

// Storage интерфейс хранилища метрик.
type Storage interface {

	// SaveMetrics сохраняет метрики.
	SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error)

	// SaveGauge сохраняет метрики типа Gauge.
	SaveGauge(ctx context.Context, name string, value float64) (err error)

	// SaveCount сохраняет метрики типа Count.
	SaveCount(ctx context.Context, name string, value int64) (err error)

	// Gauges возвращает список метрик типа Gauge.
	Gauges(ctx context.Context) (gauges []models.Gauge, err error)

	// Counters возвращает список метрик типа Counters.
	Counters(ctx context.Context) (counters []models.Counter, err error)

	// Gauge возвращает метрику типа Gauge.
	Gauge(ctx context.Context, name string) (gauge models.Gauge, err error)

	// Counter возвращает метрику типа Counter.
	Counter(ctx context.Context, name string) (counter models.Counter, err error)

	// Close закрывает соединение с хранилищем.
	Close(ctx context.Context) error

	// Ping пинг хранилища.
	Ping(ctx context.Context) error
}
