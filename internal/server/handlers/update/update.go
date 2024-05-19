package update

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/go-chi/chi"
)

type MetricsSaver interface {
	SaveGauge(name string, value float64) (err error)
	SaveCount(name string, value int64) (err error)
}

const (
	internalErrMsg = "Internal error"
)

func UpdateHandler(storage MetricsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// logger.Log.Debug("decoding request")
		var req models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			// logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if req.MType == "" || (req.MType != handlers.Gauge && req.MType != handlers.Counter) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			http.Error(w, "Invalid metric name", http.StatusNotFound)
			return
		}

		switch req.MType {
		case handlers.Gauge:
			err := storage.SaveGauge(req.ID, *req.Value)
			if err != nil {
				log.Printf("failed to save gauge: %v", err) // переделать лог
				http.Error(w, internalErrMsg, http.StatusInternalServerError)
				return
			}
		case handlers.Counter:
			err := storage.SaveCount(req.ID, *req.Delta)
			if err != nil {
				log.Printf("failed to save counter: %v", err) // переделать лог
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
		}

		resp := models.Metrics{
			ID:    req.ID,
			Value: req.Value,
			Delta: req.Delta,
			MType: req.MType,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			log.Printf("error encoding response %v", err)
			return
		}
	}
}

func UpdateHandlerRouteParams(storage MetricsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")

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
					log.Printf("failed to save gauge: %v", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
		}

		if mType == handlers.Counter {
			if val, err := strconv.ParseInt(mVal, 0, 64); err == nil {
				err := storage.SaveCount(mName, val)
				if err != nil {
					log.Printf("failed to save counter: %v", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
		}
	}
}
