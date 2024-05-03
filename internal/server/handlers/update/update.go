package update

import (
	"net/http"
	"strconv"

	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/go-chi/chi"
)

type MetricsSaver interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
}

func UpdateHandler(storage MetricsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		mType := chi.URLParam(r, "type")
		mName := chi.URLParam(r, "name")
		mVal := chi.URLParam(r, "value")

		if mType == "" || (mType != handlers.Gauge && mType != handlers.Counter) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if mName == "" {
			http.Error(w, "Invalid metric name", http.StatusNotFound)
			return
		}

		if mType == handlers.Gauge {
			if val, err := strconv.ParseFloat(mVal, 64); err == nil {
				err := storage.SaveGauge(mName, val)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
			}
		}

		if mType == handlers.Counter {
			if val, err := strconv.ParseInt(mVal, 0, 64); err == nil {
				err := storage.SaveCount(mName, val)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
			}
		}
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	}
}
