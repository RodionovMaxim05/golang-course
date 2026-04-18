package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-stat/collector/internal/domain"
	collectorpb "repo-stat/proto/collector"
)

type GetRepoUsecase interface {
	Execute(ctx context.Context, owner, name string) (domain.Repository, error)
}

type GetSubscriptionsInfoUsecase interface {
	Execute(ctx context.Context) ([]domain.Repository, error)
}

type RepoHandler struct {
	log *slog.Logger
	collectorpb.UnimplementedRepoServiceServer
	repoUsecase     GetRepoUsecase
	subsInfoUsecase GetSubscriptionsInfoUsecase
}

func NewRepoHandler(log *slog.Logger, repoUsecase GetRepoUsecase, subscriptionsInfoUsecase GetSubscriptionsInfoUsecase) *RepoHandler {
	return &RepoHandler{
		log:             log,
		repoUsecase:     repoUsecase,
		subsInfoUsecase: subscriptionsInfoUsecase,
	}
}

func (rh *RepoHandler) GetRepo(ctx context.Context, req *collectorpb.GetRepoRequest) (*collectorpb.GetRepoResponse, error) {
	rh.log.Debug("processor get repo request received", "name", req.Name, "repo", req.Repo)

	repo, err := rh.repoUsecase.Execute(ctx, req.Name, req.Repo)
	if err != nil {
		return nil, mapError(err)
	}

	return toProtoRepo(repo), nil
}

func (rh *RepoHandler) GetSubscriptionsInfo(ctx context.Context, req *collectorpb.GetSubsInfoRequest) (*collectorpb.GetSubsInfoResponse, error) {
	rh.log.Debug("processor get subscriptions info request received")

	repos, err := rh.subsInfoUsecase.Execute(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	result := make([]*collectorpb.GetRepoResponse, 0, len(repos))
	for _, repo := range repos {
		result = append(result, toProtoRepo(repo))
	}

	return &collectorpb.GetSubsInfoResponse{Repositories: result}, nil
}

func toProtoRepo(repo domain.Repository) *collectorpb.GetRepoResponse {
	return &collectorpb.GetRepoResponse{
		FullName:        repo.FullName,
		Description:     repo.Description,
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       timestamppb.New(repo.CreatedAt),
	}
}
