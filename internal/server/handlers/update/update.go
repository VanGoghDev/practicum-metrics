package update

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
)

type MetricsSaver interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
}

func UpdateHandler(storage MetricsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logger.Log.Debug("decoding request")
		var req models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			// logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch req.MType {
		case handlers.Gauge:
			err := storage.SaveGauge(req.ID, *req.Value)
			if err != nil {
				log.Printf("failed to save gauge: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
		case handlers.Counter:
			err := storage.SaveCount(req.ID, *req.Delta)
			if err != nil {
				log.Printf("failed to save counter: %v", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	}
}
