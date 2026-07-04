package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type RepoGetter interface {
	GetRepo(ctx context.Context, owner, repo string) (domain.Repository, error)
}

type GetRepo struct {
	client RepoGetter
}

func NewGetRepo(client RepoGetter) *GetRepo {
	return &GetRepo{
		client: client,
	}
}

// Execute fetches repository information for the given owner/repo.
func (gr *GetRepo) Execute(ctx context.Context, owner, repo string) (domain.Repository, error) {
	return gr.client.GetRepo(ctx, owner, repo)
}
