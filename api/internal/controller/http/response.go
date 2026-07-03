package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/domain"
)

func DomainErrToHTTP(err error) int {
	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrRateLimited:
		return http.StatusTooManyRequests
	case domain.ErrUnavailable:
		return http.StatusServiceUnavailable
	case domain.ErrInvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, log *slog.Logger, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error("failed to write response", "error", err)
	}
}
