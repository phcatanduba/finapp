package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"panic", rec,
					"stack", string(debug.Stack()),
					"path", r.URL.Path,
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
