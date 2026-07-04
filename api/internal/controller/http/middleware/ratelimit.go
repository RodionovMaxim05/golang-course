package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, float64, error)
}

// RateLimit returns a middleware that enforces a per-client-IP rate limit
// via the given RateLimiter. Requests are keyed by the remote peer's IP
// address. If the limiter backend is unavailable, requests are allowed
// through (fail open) to preserve availability.
func RateLimit(limiter RateLimiter, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Warn("failed to parse remote addr, using raw value", "remote_addr", r.RemoteAddr, "error", err)
				ip = r.RemoteAddr
			}

			allowed, remaining, err := limiter.Allow(r.Context(), "rate_limit:"+ip)
			if err != nil {
				log.Error("rate limiter failed, allowing request", "error", err, "key", ip)
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", remaining))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"}); err != nil {
					log.Error("failed to write response", "error", err)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
