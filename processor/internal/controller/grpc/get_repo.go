package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-watcher/processor/internal/domain"
	processorpb "repo-watcher/proto/gen/go/processor/v1"
)

// GetRepo handles a client request for repository information. It returns
// the cached data if available and successfully processed, a PENDING
// status if the data is still being collected, or a gRPC error mapped
// from the underlying use case/storage error.
func (s *Server) GetRepo(ctx context.Context, req *processorpb.GetRepoRequest) (*processorpb.GetRepoResponse, error) {
	s.log.Debug("processor get repo request received", "owner", req.Owner, "repo", req.Repo)

	resp, err := s.getRepo.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrAccepted) {
			s.log.Debug("repository is processing, returning Accepted status", "owner", req.Owner, "repo", req.Repo)
			return &processorpb.GetRepoResponse{
				Status: processorpb.GetRepoResponse_STATUS_PENDING,
			}, nil
		}

		s.log.Error("usecase execution failed", "error", err)
		return nil, mapUseCaseError(err)
	}

	switch resp.Status {
	case domain.StatusPending:
		s.log.Debug("repository found in DB but still has PENDING status", "repo", resp.FullName)
		return &processorpb.GetRepoResponse{
			Status: processorpb.GetRepoResponse_STATUS_PENDING,
		}, nil

	case domain.StatusError:
		s.log.Warn("repository found in DB with ERROR status", "repo", resp.FullName, "code", resp.ErrorCode)
		return nil, mapRepoErrorCode(resp.FullName, resp.ErrorCode)

	case domain.StatusSuccess:
		return &processorpb.GetRepoResponse{
			Status:          processorpb.GetRepoResponse_STATUS_SUCCESS,
			FullName:        resp.FullName,
			Description:     resp.Description,
			StargazersCount: int32(resp.StargazersCount),
			ForksCount:      int32(resp.ForksCount),
			CreatedAt:       timestamppb.New(resp.CreatedAt),
		}, nil

	default:
		s.log.Error("unknown repository status in database", "status", resp.Status)
		return nil, status.Error(codes.Internal, "internal data error")
	}
}

// mapUseCaseError translates an error returned by the GetRepo use case
// into an appropriate gRPC status. Errors not recognized as one of the
// known domain sentinel errors default to codes.Internal.
func mapUseCaseError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "repository not found")
	case errors.Is(err, domain.ErrInvalidArgument):
		return status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	case errors.Is(err, domain.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, "rate limit exceeded")
	default:
		return status.Errorf(codes.Internal, "internal server error: %v", err)
	}
}

// mapRepoErrorCode translates a domain.ErrorCode, previously reported by
// the Collector and persisted for a repository in ERROR status, into an
// appropriate gRPC status for the client.
func mapRepoErrorCode(fullName string, code domain.ErrorCode) error {
	switch code {
	case domain.ErrorCodeRepositoryNotFound:
		return status.Errorf(codes.NotFound, "repository %s not found on github", fullName)
	case domain.ErrorCodeGitHubRateLimitExceeded:
		return status.Error(codes.ResourceExhausted, "github rate limit exceeded while collecting repository data")
	case domain.ErrorCodeInternalCollectorError, domain.ErrorCodeUnspecified:
		return status.Errorf(codes.Internal, "failed to collect repository data for %s", fullName)
	default:
		return status.Errorf(codes.Internal, "failed to collect repository data: unknown error code %q", code)
	}
}
