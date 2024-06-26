package update

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const (
	internalErrMsg = "Internal error"
)

func UpdateHandler(zlog *zap.SugaredLogger, storage routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &models.Metrics{}

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			zlog.Warnf("failed to decode JSON body", zap.Error(err))
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
			err := storage.SaveGauge(r.Context(), req.ID, *req.Value)
			if err != nil {
				zlog.Warnf("failed to save gauge: %v", err) // переделать лог
				http.Error(w, internalErrMsg, http.StatusInternalServerError)
				return
			}
		case handlers.Counter:
			err := storage.SaveCount(r.Context(), req.ID, *req.Delta)
			if err != nil {
				zlog.Warnf("failed to save counter: %v", err) // переделать лог
				http.Error(w, internalErrMsg, http.StatusInternalServerError)
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
			zlog.Warnf("error encoding response %v", err)
			return
		}
	}
}

func UpdateHandlerRouteParams(zlog *zap.SugaredLogger, storage routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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
				err := storage.SaveGauge(r.Context(), mName, val)
				if err != nil {
					zlog.Warnf("failed to save gauge: %v", err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
		}

		if mType == handlers.Counter {
			if val, err := strconv.ParseInt(mVal, 0, 64); err == nil {
				err := storage.SaveCount(r.Context(), mName, val)
				if err != nil {
					zlog.Warnf("failed to save counter: %v", err)
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
