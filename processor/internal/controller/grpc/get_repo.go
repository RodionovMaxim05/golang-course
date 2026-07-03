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

func (s *Server) GetRepo(ctx context.Context, req *processorpb.GetRepoRequest) (*processorpb.GetRepoResponse, error) {
	s.log.Debug("processor get repo request received", "owner", req.Owner, "repo", req.Repo)

	resp, err := s.getRepo.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrAccepted) {
			s.log.Debug("repository is processing, returning Accepted status", "owner", req.Owner, "repo", req.Repo)
			return nil, status.Error(codes.Unavailable, "repository data is being collected, try again later")
		}

		s.log.Error("usecase execution failed", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	switch resp.Status {
	case "PENDING":
		s.log.Debug("repository found in DB but still has PENDING status", "repo", resp.FullName)
		return nil, status.Error(codes.Unavailable, "repository data is still being collected")

	case "ERROR":
		s.log.Warn("repository found in DB with ERROR status", "repo", resp.FullName, "code", resp.ErrorCode)
		if resp.ErrorCode == "REPOSITORY_NOT_FOUND" {
			return nil, status.Errorf(codes.NotFound, "repository %s not found on github", resp.FullName)
		}
		return nil, status.Errorf(codes.ResourceExhausted, "failed to collect repository data: %s", resp.ErrorCode)

	case "SUCCESS":
		return &processorpb.GetRepoResponse{
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
