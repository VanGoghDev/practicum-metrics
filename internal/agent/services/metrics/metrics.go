package metrics

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"go.uber.org/zap"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
)

type MetricsProvider struct {
	log     *zap.Logger
	metrics []*models.Metrics
}

func New(log *zap.Logger) *MetricsProvider {
	metrics := []*models.Metrics{
		createMetric("Alloc", gaugeType),
		createMetric("BuckHashSys", gaugeType),
		createMetric("Frees", gaugeType),
		createMetric("GCCPUFraction", gaugeType),
		createMetric("GCSys", gaugeType),
		createMetric("HeapAlloc", gaugeType),
		createMetric("HeapIdle", gaugeType),
		createMetric("HeapInuse", gaugeType),
		createMetric("HeapObjects", gaugeType),
		createMetric("HeapReleased", gaugeType),
		createMetric("HeapSys", gaugeType),
		createMetric("LastGC", gaugeType),
		createMetric("Lookups", gaugeType),
		createMetric("MCacheInuse", gaugeType),
		createMetric("MCacheSys", gaugeType),
		createMetric("MSpanSys", gaugeType),
		createMetric("MSpanInuse", gaugeType),
		createMetric("Mallocs", gaugeType),
		createMetric("NextGC", gaugeType),
		createMetric("NumForcedGC", gaugeType),
		createMetric("NumGC", gaugeType),
		createMetric("OtherSys", gaugeType),
		createMetric("PauseTotalNs", gaugeType),
		createMetric("StackInuse", gaugeType),
		createMetric("StackSys", gaugeType),
		createMetric("Sys", gaugeType),
		createMetric("TotalAlloc", gaugeType),
		createMetric("RandomValue", gaugeType),
		createMetric("PollCount", counterType),
	}

	return &MetricsProvider{
		log:     log,
		metrics: metrics,
	}
}

func (mp *MetricsProvider) ReadMetrics(pollCount int64) ([]*models.Metrics, error) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	for _, v := range mp.metrics {
		err := populateMetric(v, m, pollCount)
		if err != nil {
			mp.log.Warn(fmt.Sprintf("Failed to fill metric with value: %v", err))
			continue
		}
	}
	return mp.metrics, nil
}

// Создает метрику с ID = name, MType = mType. Значения метрики будут заданы по умолчанию.
func createMetric(name string, mType string) *models.Metrics {
	switch mType {
	case gaugeType:
		v := float64(0)
		return createMetricVal(name, mType, &v)
	case counterType:
		v := int64(0)
		return createMetricDelta(name, mType, &v)
	}
	return &models.Metrics{}
}

// Создает метрику с ID = name, MType = mType, Value = value. Значение Delta будет задано по умолчанию.
func createMetricVal(name string, mType string, value *float64) *models.Metrics {
	return &models.Metrics{
		ID:    name,
		MType: mType,
		Value: value,
	}
}

// Создает метрику с ID = name, MType = mType, Delta = delta. Значение Value будет задано по умолчанию.
func createMetricDelta(name string, mType string, delta *int64) *models.Metrics {
	return &models.Metrics{
		ID:    name,
		MType: mType,
		Delta: delta,
	}
}

// Наполняет метрику значением из runtime в зависимости от ID.
func populateMetric(metric *models.Metrics, m *runtime.MemStats, pollCount int64) error {
	switch metric.ID {
	case "Alloc":
		val := float64(m.Alloc)
		metric.Value = &val
	case "BuckHashSys":
		val := float64(m.BuckHashSys)
		metric.Value = &val
	case "Frees":
		val := float64(m.Frees)
		metric.Value = &val
	case "GCCPUFraction":
		val := float64(m.GCCPUFraction)
		metric.Value = &val
	case "GCSys":
		val := float64(m.GCSys)
		metric.Value = &val
	case "HeapAlloc":
		val := float64(m.HeapAlloc)
		metric.Value = &val
	case "HeapSys":
		val := float64(m.HeapSys)
		metric.Value = &val
	case "HeapIdle":
		val := float64(m.HeapIdle)
		metric.Value = &val
	case "HeapInuse":
		val := float64(m.HeapInuse)
		metric.Value = &val
	case "HeapObjects":
		val := float64(m.HeapObjects)
		metric.Value = &val
	case "HeapReleased":
		val := float64(m.HeapReleased)
		metric.Value = &val
	case "LastGC":
		val := float64(m.LastGC)
		metric.Value = &val
	case "Lookups":
		val := float64(m.Lookups)
		metric.Value = &val
	case "MCacheInuse":
		val := float64(m.MCacheInuse)
		metric.Value = &val
	case "MCacheSys":
		val := float64(m.MCacheSys)
		metric.Value = &val
	case "MSpanInuse":
		val := float64(m.MSpanInuse)
		metric.Value = &val
	case "MSpanSys":
		val := float64(m.MSpanSys)
		metric.Value = &val
	case "Mallocs":
		val := float64(m.Mallocs)
		metric.Value = &val
	case "NextGC":
		val := float64(m.NextGC)
		metric.Value = &val
	case "NumForcedGC":
		val := float64(m.NumForcedGC)
		metric.Value = &val
	case "NumGC":
		val := float64(m.NumGC)
		metric.Value = &val
	case "OtherSys":
		val := float64(m.OtherSys)
		metric.Value = &val
	case "PauseTotalNs":
		val := float64(m.PauseTotalNs)
		metric.Value = &val
	case "StackInuse":
		val := float64(m.StackInuse)
		metric.Value = &val
	case "StackSys":
		val := float64(m.StackSys)
		metric.Value = &val
	case "Sys":
		val := float64(m.Sys)
		metric.Value = &val
	case "TotalAlloc":
		val := float64(m.TotalAlloc)
		metric.Value = &val
	case "RandomValue":
		val := rand.Float64()
		metric.Value = &val
	case "PollCount":
		metric.Delta = &pollCount
	default:
		return fmt.Errorf("unknown metric: %s", metric.ID)
	}
	return nil
}
