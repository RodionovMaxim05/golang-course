package usecase

import (
	"context"

	"repo-watcher/collector/internal/domain"
)

type RepoService interface {
	GetRepo(ctx context.Context, owner, repo string) (domain.Repository, error)
}

type GetRepoUsecase struct {
	client RepoService
}

func NewRepoUsecase(client RepoService) *GetRepoUsecase {
	return &GetRepoUsecase{client: client}
}

func (gru *GetRepoUsecase) Execute(ctx context.Context, owner, repo string) (domain.Repository, error) {
	return gru.client.GetRepo(ctx, owner, repo)
}
