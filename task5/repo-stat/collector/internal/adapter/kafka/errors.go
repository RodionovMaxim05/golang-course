package kafka

import (
	"errors"

	"repo-stat/collector/internal/domain"
)

func mapDomainErrorToCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "REPOSITORY_NOT_FOUND"
	case errors.Is(err, domain.ErrRateLimited):
		return "GITHUB_RATE_LIMIT_EXCEEDED"
	default:
		return "INTERNAL_COLLECTOR_ERROR"
	}
}
