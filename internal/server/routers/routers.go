package routers

import "github.com/VanGoghDev/practicum-metrics/internal/domain/models"

type Storage interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
	Gauges() (gauges []models.Gauge, err error)
	Counters() (counters []models.Counter, err error)
	Gauge(name string) (gauge models.Gauge, err error)
	Counter(name string) (counter models.Counter, err error)
	Close() error
}
