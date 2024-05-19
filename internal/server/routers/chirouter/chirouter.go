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
		r.Post("/", metrics.MetricHandler(p))
		r.Get("/{type}/{name}", metrics.MetricHandlerRouterParams(p))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", update.UpdateHandler(s))
		r.Post("/{type}/{name}/{value}", update.UpdateHandlerRouteParams(s))
	})

	return r
}
