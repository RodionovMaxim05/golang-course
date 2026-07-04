package middleware

import (
	"log/slog"
	"net/http"
)

// Recover returns a middleware that recovers from panics in the wrapped
// handler chain, logging the panic value and request path, and responding
// with a 500 Internal Server Error instead of letting the panic crash the
// connection.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered in http handler", "panic", rec, "path", r.URL.Path)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
