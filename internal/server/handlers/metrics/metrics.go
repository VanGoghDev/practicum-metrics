package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"
	"github.com/VanGoghDev/practicum-metrics/internal/util/converter"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var (
	errFailedToFetchGauge   = errors.New("failed to fetch gauge")
	errFailedToFetchCounter = errors.New("failed to fetch counter")
)

const (
	internalErrMsg = "Internal error"
	notFoundErrMsg = "Not found"
)

func MetricsHandler(zlog *zap.SugaredLogger, s routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		gauges, err := s.Gauges()
		if err != nil {
			zlog.Warnf("failed to fetch gauges: %v", err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}

		counters, err := s.Counters()
		if err != nil {
			zlog.Warnf("failed to fetch counters: %v", err)
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}

		for _, g := range gauges {
			sV, err := converter.Str(g.Value)
			if err != nil {
				zlog.Warnf("failed to convert gauge value to string: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				break
			}

			_, err = fmt.Fprintf(w, "%s: %s \n", g.Name, sV)
			if err != nil {
				zlog.Warnf("failed to print gauges: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				break
			}
		}

		for _, c := range counters {
			sV, err := converter.Str(c.Value)
			if err != nil {
				zlog.Warnf("failed to convert counter value to string: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				break
			}

			_, err = fmt.Fprintf(w, "%s: %s \n", c.Name, sV)
			if err != nil {
				zlog.Warnf("failed to print counters: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				break
			}
		}
	}
}

func MetricHandler(zlog *zap.SugaredLogger, s routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			zlog.Errorf("failed to decode request: %w", err)
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
		case handlers.Counter:
			{
				counter, err := s.Counter(req.ID)
				if err != nil {
					handleError(zlog, err, w)
					return
				}
				resp := models.Metrics{
					ID:    req.ID,
					Delta: &counter.Value,
					MType: req.MType,
				}
				enc := json.NewEncoder(w)
				if err := enc.Encode(resp); err != nil {
					zlog.Errorf("error encoding response: %w", err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				return
			}
		case handlers.Gauge:
			{
				gauge, err := s.Gauge(req.ID)
				if err != nil {
					handleError(zlog, err, w)
					return
				}

				resp := models.Metrics{
					ID:    req.ID,
					Value: &gauge.Value,
					MType: req.MType,
				}
				enc := json.NewEncoder(w)
				if err := enc.Encode(resp); err != nil {
					zlog.Errorf("error encoding writer: %w", errFailedToFetchGauge)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				return
			}
		}
	}
}

func MetricHandlerRouterParams(zlog *zap.SugaredLogger, s routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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
					if errors.Is(err, serrors.ErrNotFound) {
						http.Error(w, "Not found", http.StatusNotFound)
						return
					}
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				sV, err := converter.Str(counter.Value)
				if err != nil {
					zlog.Errorf("failed to convert counter value to string: %w", err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}

				_, err = fmt.Fprintf(w, "%s", sV)
				if err != nil {
					zlog.Errorf("%v: %v", errFailedToFetchCounter, err)

					log.Printf("%v: %v", errFailedToFetchCounter, err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				return
			}
		case handlers.Gauge:
			{
				gauge, err := s.Gauge(mName)
				if err != nil {
					if errors.Is(err, serrors.ErrNotFound) {
						http.Error(w, "Not found", http.StatusNotFound)
						return
					}
					zlog.Errorf("%v: %w", err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				sV, err := converter.Str(gauge.Value)
				if err != nil {
					zlog.Errorf("failed to convert gauge value to string: %w", err)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				_, err = fmt.Fprintf(w, "%s", sV)
				if err != nil {
					zlog.Errorf("%v", errFailedToFetchGauge)
					http.Error(w, internalErrMsg, http.StatusInternalServerError)
					return
				}
				return
			}
		}
	}
}

func handleError(zlog *zap.SugaredLogger, err error, w http.ResponseWriter) {
	if errors.Is(err, serrors.ErrNotFound) {
		http.Error(w, notFoundErrMsg, http.StatusNotFound)
		return
	}
	zlog.Errorf("Invalid metric type: %w", err)
	http.Error(w, internalErrMsg, http.StatusInternalServerError)
}
