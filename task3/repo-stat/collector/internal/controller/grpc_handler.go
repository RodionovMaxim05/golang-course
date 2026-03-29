package controller

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-stat/collector/internal/domain"
	collectorpb "repo-stat/proto/collector"
)

type GetRepoUsecase interface {
	Execute(ctx context.Context, owner, name string) (domain.Repository, error)
}

type RepoHandler struct {
	log *slog.Logger
	collectorpb.UnimplementedRepoServiceServer
	repoUsecase GetRepoUsecase
}

func NewRepoHandler(log *slog.Logger, repoUsecase GetRepoUsecase) *RepoHandler {
	return &RepoHandler{log: log, repoUsecase: repoUsecase}
}

func (rh *RepoHandler) GetRepo(ctx context.Context, req *collectorpb.GetRepoRequest) (*collectorpb.GetRepoResponse, error) {
	rh.log.Debug("processor get repo request received", "name", req.Name, "repo", req.Repo)

	repo, err := rh.repoUsecase.Execute(ctx, req.Name, req.Repo)
	if err != nil {
		return nil, mapError(err)
	}

	return &collectorpb.GetRepoResponse{
		Name:            req.Name,
		Repo:            req.Repo,
		Description:     repo.Description,
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       timestamppb.New(repo.CreatedAt),
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
