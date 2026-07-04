package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/domain"
)

// DomainErrToHTTP translates a domain sentinel error into the
// corresponding HTTP status code. Unrecognized errors default to 500.
func DomainErrToHTTP(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrPending):
		return http.StatusAccepted
	case errors.Is(err, domain.ErrRateLimited):
		return http.StatusTooManyRequests
	case errors.Is(err, domain.ErrUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// writeJSON writes payload as a JSON response body with the given status
// code, logging any encoding failure.
func writeJSON(w http.ResponseWriter, log *slog.Logger, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error("failed to write response", "error", err)
	}
}

func writeError(w http.ResponseWriter, log *slog.Logger, err error) {
	code := DomainErrToHTTP(err)
	msg := err.Error()
	if code == http.StatusInternalServerError {
		msg = "internal server error"
	}
	writeJSON(w, log, code, map[string]string{"error": msg})
}
