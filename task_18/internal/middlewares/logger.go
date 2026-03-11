package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("[Middleware Logger]",
			"method", r.Method,
			"path", r.URL.Path,
			"time", time.Now().String())

		fn(w, r)
	}
}
