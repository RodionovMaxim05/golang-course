package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-watcher/subscriber/internal/domain"
)

// grpcError translates a domain error returned by a use case into the
// appropriate gRPC status code. Errors not recognized as one of the known
// domain sentinel errors default to codes.Internal.
func grpcError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrRepoNotFound), errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadySubscribed):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
