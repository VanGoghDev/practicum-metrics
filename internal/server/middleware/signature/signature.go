package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"go.uber.org/zap"
)

func New(zlog *zap.SugaredLogger, cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if cfg.Key == "" {
				next.ServeHTTP(w, r)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				zlog.Warnf("failed to read request body: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = r.Body.Close()
			if err != nil {
				zlog.Warnf("failed to close request body: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			rSig := r.Header.Get("HashSHA256")
			// rs, _ := hex.DecodeString(rSig)
			h := hmac.New(sha256.New, []byte(cfg.Key))
			h.Write(body)
			dst := h.Sum(nil)
			if !hmac.Equal(dst, []byte(rSig)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("HashSHA256", string(dst))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
