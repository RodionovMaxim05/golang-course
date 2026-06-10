package http

import (
	"net/http"
	"repo-stat/api/internal/domain"
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
