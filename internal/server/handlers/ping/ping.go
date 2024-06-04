package ping

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"go.uber.org/zap"
)

func PingHandler(zlog *zap.SugaredLogger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("pgx", cfg.DBConnectionString)
		if err != nil {
			zlog.Warnf("failed to open db: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = db.Ping()
		if err != nil {
			zlog.Warnf("failed to ping db: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
