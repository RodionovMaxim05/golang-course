package usecase

import (
	"context"
	"fmt"

	"repo-watcher/collector/internal/domain"
)

type RepoFetcher interface {
	GetRepo(ctx context.Context, owner, repo string) (domain.Repository, error)
}

type GetRepoUsecase struct {
	client RepoFetcher
}

func NewGetRepoUsecase(client RepoFetcher) *GetRepoUsecase {
	return &GetRepoUsecase{client: client}
}

// Execute fetches and returns repository data for the given owner/repo.
func (gru *GetRepoUsecase) Execute(ctx context.Context, owner, repo string) (domain.Repository, error) {
	repository, err := gru.client.GetRepo(ctx, owner, repo)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("get repo %s/%s: %w", owner, repo, err)
	}
	return repository, nil
}
