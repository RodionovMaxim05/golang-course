package usecase

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type RepoGetter interface {
	GetRepo(ctx context.Context, name, repo string) (*domain.Repository, error)
}

type GetRepo struct {
	repoGetter RepoGetter
}

func NewGetRepo(repoGetter RepoGetter) *GetRepo {
	return &GetRepo{
		repoGetter: repoGetter,
	}
}

func (gr *GetRepo) Execute(ctx context.Context, name, repo string) (*domain.Repository, error) {
	return gr.repoGetter.GetRepo(ctx, name, repo)
}
