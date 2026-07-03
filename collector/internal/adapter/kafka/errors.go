package kafka

import (
	"errors"

	"repo-watcher/collector/internal/domain"
	commonpb "repo-watcher/proto/gen/go/common/v1"
)

func mapDomainErrorToCode(err error) commonpb.ErrorCode {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return commonpb.ErrorCode_ERROR_CODE_REPOSITORY_NOT_FOUND
	case errors.Is(err, domain.ErrRateLimited):
		return commonpb.ErrorCode_ERROR_CODE_GITHUB_RATE_LIMIT_EXCEEDED
	default:
		return commonpb.ErrorCode_ERROR_CODE_INTERNAL_COLLECTOR_ERROR
	}
}
