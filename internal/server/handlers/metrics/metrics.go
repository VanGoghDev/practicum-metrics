package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/go-chi/chi"
)

type MetricsProvider interface {
	Gauges() (gauges []models.Gauge, err error)
	Counters() (counters []models.Counter, err error)
	Gauge(name string) (gauge models.Gauge, err error)
	Counter(name string) (counter models.Counter, err error)
}

func MetricsHandler(s MetricsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gauges, err := s.Gauges()

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		counters, err := s.Counters()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		for _, g := range gauges {
			_, err = fmt.Fprintf(w, "%s: %f \n", g.Name, g.Value)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				break
			}
		}

		for _, c := range counters {
			_, err = fmt.Fprintf(w, "%s: %d \n", c.Name, c.Value)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				break
			}
		}
	}
}

func MetricHandler(s MetricsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
		mType := chi.URLParam(r, "type")
		mName := chi.URLParam(r, "name")

		if mType == "" || (mType != handlers.Gauge && mType != handlers.Counter) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if mName == "" {
			http.Error(w, "Invalid metric name", http.StatusNotFound)
			return
		}

		switch mType {
		case handlers.Counter:
			{
				counter, err := s.Counter(mName)
				if err != nil {
					if errors.Is(err, memstorage.ErrNotFound) {
						http.Error(w, "Not found", http.StatusNotFound)
						return
					}
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				_, err = fmt.Fprintf(w, "%d", counter.Value)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				return
			}
		case handlers.Gauge:
			{
				gauge, err := s.Gauge(mName)
				if err != nil {
					if errors.Is(err, memstorage.ErrNotFound) {
						http.Error(w, "Not found", http.StatusNotFound)
						return
					}
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				_, err = fmt.Fprintf(w, "%s", strconv.FormatFloat(gauge.Value, 'f', -1, 64))
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				return
			}
		}
	}
}