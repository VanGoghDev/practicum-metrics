package chirouter

import (
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
	mwLogger "github.com/VanGoghDev/practicum-metrics/internal/server/middleware/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func BuildRouter(s update.MetricsSaver, p metrics.MetricsProvider, log *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(mwLogger.New(log))

	r.Route("/", func(r chi.Router) {
		r.Get("/", metrics.MetricsHandler(p))
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", metrics.MetricHandler(p))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", update.UpdateHandler(s))
	})

	return r
}
