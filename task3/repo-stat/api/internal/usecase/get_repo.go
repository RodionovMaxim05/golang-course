package usecase

import (
	"context"

	processorpb "repo-stat/proto/processor"
)

type RepoGetter interface {
	GetRepo(ctx context.Context, name, repo string) (*processorpb.GetRepoResponse, error)
}

type GetRepo struct {
	repoGetter RepoGetter
}

func NewGetRepo(repoGetter RepoGetter) *GetRepo {
	return &GetRepo{
		repoGetter: repoGetter,
	}
}

func (gr *GetRepo) Execute(ctx context.Context, name, repo string) (*processorpb.GetRepoResponse, error) {
	return gr.repoGetter.GetRepo(ctx, name, repo)
}
