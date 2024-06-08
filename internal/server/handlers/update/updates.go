package update

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"go.uber.org/zap"
)

func UpdatesHandler(ctx context.Context, zlog *zap.SugaredLogger, storage routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := []*models.Metrics{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			zlog.Warnf("failed to decode JSON body", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		err := storage.SaveMetrics(ctx, metrics)
		if err != nil {
			zlog.Warnf("failed to save metrics: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}
}
