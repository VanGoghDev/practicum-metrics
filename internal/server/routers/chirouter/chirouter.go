package chirouter

import (
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
	"github.com/VanGoghDev/practicum-metrics/internal/server/middleware/compressor"
	"github.com/VanGoghDev/practicum-metrics/internal/server/middleware/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func BuildRouter(s routers.Storage, log *zap.Logger) chi.Router {
	r := chi.NewRouter()
	sugarlog := log.Sugar()
	r.Use(logger.New(sugarlog))
	r.Use(compressor.New(sugarlog))

	r.Route("/", func(r chi.Router) {
		r.Get("/", metrics.MetricsHandler(sugarlog, s))
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", metrics.MetricHandler(sugarlog, s))
		r.Get("/{type}/{name}", metrics.MetricHandlerRouterParams(sugarlog, s))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", update.UpdateHandler(sugarlog, s))
		r.Post("/{type}/{name}/{value}", update.UpdateHandlerRouteParams(sugarlog, s))
	})

	return r
}
