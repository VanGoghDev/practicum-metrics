package chirouter

import (
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/provider"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
	"github.com/go-chi/chi"
)

func BuildRouter(s update.MetricsSaver, p provider.MetricsProvider) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", provider.MetricsHandler(p))
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", provider.MetricHandler(p))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", update.UpdateHandler(s))
	})

	// r.HandleFunc(`/update/{type}/{name}/{value}`, update.UpdateHandler(s))

	return r
}
