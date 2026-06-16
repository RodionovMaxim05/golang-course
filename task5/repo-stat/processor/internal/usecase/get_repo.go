package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"repo-stat/processor/internal/domain"
)

type RepoGetter interface {
	SendRepoRequest(ctx context.Context, owner, repo string) error
}

type GetRepo struct {
	log         *slog.Logger
	repoGetter  RepoGetter
	dataStorage domain.DataStorage
}

func NewGetRepo(repoGetter RepoGetter, storage domain.DataStorage, log *slog.Logger) *GetRepo {
	return &GetRepo{
		log:         log,
		repoGetter:  repoGetter,
		dataStorage: storage,
	}
}

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
			Status:   "PENDING",
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
