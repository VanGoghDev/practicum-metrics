package chirouter

import (
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/metrics"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
	"github.com/go-chi/chi"
)

func BuildRouter(s update.MetricsSaver, p metrics.MetricsProvider) chi.Router {
	r := chi.NewRouter()

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
