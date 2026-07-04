package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-watcher/api/internal/domain"
)

// ErrToDomain translates a gRPC error returned by a downstream service
// into the corresponding domain sentinel error, based on its gRPC status
// code. Errors without a recognized status code, or without a gRPC
// status at all, are returned unchanged.
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
	case codes.AlreadyExists:
		return domain.ErrAlreadyExists
	default:
		return err
	}
}
