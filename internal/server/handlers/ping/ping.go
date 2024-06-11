package ping

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers"
	"go.uber.org/zap"
)

func PingHandler(zlog *zap.SugaredLogger, cfg *config.Config, storage routers.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := storage.Ping(r.Context())
		if err != nil {
			zlog.Warnf("failed to ping db: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
