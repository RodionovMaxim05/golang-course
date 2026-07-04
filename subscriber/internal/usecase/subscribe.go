package usecase

import (
	"context"
	"errors"
	"fmt"

	"repo-watcher/subscriber/internal/domain"
)

type GitHubChecker interface {
	RepoExists(ctx context.Context, owner, repo string) (bool, error)
}

type SubscriptionInteractor interface {
	Create(ctx context.Context, sub *domain.Subscription) (*domain.SubscriptionRecord, error)
}

type Subscribe struct {
	githubClient GitHubChecker
	repo         SubscriptionInteractor
}

func NewSubscribe(githubClient GitHubChecker, repository SubscriptionInteractor) *Subscribe {
	return &Subscribe{githubClient: githubClient, repo: repository}
}

// Execute verifies that owner/repo exists on GitHub and creates a new
// subscription for it. Returns domain.ErrRepoNotFound if the repository
// doesn't exist on GitHub, or domain.ErrAlreadySubscribed if a
// subscription already exists.
func (s *Subscribe) Execute(ctx context.Context, owner, repo string) (*domain.SubscriptionRecord, error) {
	// Check if repository exists on GitHub
	exists, err := s.githubClient.RepoExists(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("check exists: %w", err)
	}
	if !exists {
		return nil, domain.ErrRepoNotFound
	}

	sub := &domain.Subscription{Owner: owner, Repo: repo}

	resp, err := s.repo.Create(ctx, sub)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadySubscribed) {
			return nil, err
		}
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return resp, nil
}
