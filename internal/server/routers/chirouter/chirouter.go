package chirouter

import (
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/ping"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
	"github.com/VanGoghDev/practicum-metrics/internal/server/middleware/compressor"
	"github.com/VanGoghDev/practicum-metrics/internal/server/middleware/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/middleware/signature"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func BuildRouter(s routers.Storage, log *zap.Logger, cfg *config.Config) chi.Router {
	r := chi.NewRouter()
	sugarlog := log.Sugar()
	r.Use(logger.New(sugarlog))
	r.Use(signature.New(sugarlog, cfg))
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

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", update.UpdatesHandler(sugarlog, s))
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", ping.PingHandler(sugarlog, cfg, s))
	})

	return r
}
