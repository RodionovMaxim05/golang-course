package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"repo-watcher/processor/internal/domain"
)

type RepoFetchRequester interface {
	SendRepoRequest(ctx context.Context, owner, repo string) error
}

type RepoStorage interface {
	GetRepo(ctx context.Context, fullName string) (*domain.Repository, error)
	UpdateRepoStatus(ctx context.Context, repo *domain.Repository) error
}

type GetRepo struct {
	log         *slog.Logger
	repoGetter  RepoFetchRequester
	dataStorage RepoStorage
}

func NewGetRepo(repoGetter RepoFetchRequester, storage RepoStorage, log *slog.Logger) *GetRepo {
	return &GetRepo{
		log:         log,
		repoGetter:  repoGetter,
		dataStorage: storage,
	}
}

// Execute returns cached information for the repository identified by
// owner/repo. If the repository is not present in local storage, it marks
// the repository as PENDING, publishes an asynchronous fetch request via
// the RepoFetchRequester port, and returns domain.ErrAccepted to indicate the
// caller should retry later once the data has been collected.
func (gr *GetRepo) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	fullName := fmt.Sprintf("%s/%s", owner, repo)

	repoInfo, err := gr.dataStorage.GetRepo(ctx, fullName)
	if err == nil {
		gr.log.Debug("fetching repository from database", "owner", owner, "repo", repo)
		return repoInfo, nil
	}

	if errors.Is(err, domain.ErrNotFound) {
		pendingRepo := &domain.Repository{
			FullName: fullName,
			Status:   domain.StatusPending,
		}
		if dbErr := gr.dataStorage.UpdateRepoStatus(ctx, pendingRepo); dbErr != nil {
			gr.log.Error("failed to save PENDING status to DB", "repo", fullName, "error", dbErr)
		} else {
			gr.log.Debug("successfully saved PENDING status to DB", "repo", fullName)
		}

		err := gr.repoGetter.SendRepoRequest(ctx, owner, repo)
		if err != nil {
			gr.log.Error("failed to send repo request", "repo", fullName, "error", err)
			return nil, fmt.Errorf("send repo request: %w", err)
		}
		gr.log.Debug("send repo info request", "owner", owner, "repo", repo)

		return nil, domain.ErrAccepted
	}

	return nil, err
}
