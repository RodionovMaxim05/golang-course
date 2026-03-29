package usecase

import (
	"context"

	collectorpb "repo-stat/proto/collector"
)

type RepoGetter interface {
	GetRepo(ctx context.Context, name, repo string) (*collectorpb.GetRepoResponse, error)
}

type GetRepo struct {
	repoGetter RepoGetter
}

func NewGetRepo(repoGetter RepoGetter) *GetRepo {
	return &GetRepo{
		repoGetter: repoGetter,
	}
}

func (gr *GetRepo) Execute(ctx context.Context, name, repo string) (*collectorpb.GetRepoResponse, error) {
	return gr.repoGetter.GetRepo(ctx, name, repo)
}
