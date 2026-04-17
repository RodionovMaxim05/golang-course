package grpc

import (
	"repo-stat/api/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrToDomain(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return domain.ErrNotFound
	case codes.ResourceExhausted:
		return domain.ErrRateLimited
	case codes.Unavailable:
		return domain.ErrUnavailable
	case codes.InvalidArgument:
		return domain.ErrInvalidArgument
	default:
		return err
	}
}
