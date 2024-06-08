package routers

import (
	"context"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

type Storage interface {
	SaveMetrics(ctx context.Context, metrics []*models.Metrics) (err error)
	SaveGauge(ctx context.Context, name string, value float64) (err error)
	SaveCount(ctx context.Context, name string, value int64) (err error)
	Gauges(ctx context.Context) (gauges []models.Gauge, err error)
	Counters(ctx context.Context) (counters []models.Counter, err error)
	Gauge(ctx context.Context, name string) (gauge models.Gauge, err error)
	Counter(ctx context.Context, name string) (counter models.Counter, err error)
	Close(ctx context.Context) error
}
