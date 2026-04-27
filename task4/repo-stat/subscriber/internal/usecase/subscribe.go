package usecase

import (
	"context"
	"fmt"

	"repo-stat/subscriber/internal/domain"
)

type GitHubChecker interface {
	RepoExists(ctx context.Context, owner, repo string) (bool, error)
}

type Subscribe struct {
	githubClient GitHubChecker
	repo         domain.SubscriptionRepository
}

func NewSubscriptionService(githubClient GitHubChecker, repository domain.SubscriptionRepository) *Subscribe {
	return &Subscribe{githubClient: githubClient, repo: repository}
}

func (s *Subscribe) Execute(ctx context.Context, owner, repo string) (*domain.SubscriptionResponse, error) {
	// Check if repository exists on GitHub
	exists, err := s.githubClient.RepoExists(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("check exists: %w", err)
	}
	if !exists {
		return nil, domain.ErrRepoNotFound
	}

	// Check if already subscribed
	alreadySubscribed, err := s.repo.Exists(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("check subscription: %w", err)
	}
	if alreadySubscribed {
		return nil, domain.ErrAlreadySubscribed
	}

	subscription := &domain.Subscription{Owner: owner, Repo: repo}

	return s.repo.Create(ctx, subscription)
}
