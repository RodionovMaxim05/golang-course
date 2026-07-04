package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type SubscriptionCreator interface {
	Subscribe(ctx context.Context, owner, repo string) (*domain.Subscription, error)
}

type Subscribe struct {
	client SubscriptionCreator
}

func NewSubscriber(client SubscriptionCreator) *Subscribe {
	return &Subscribe{
		client: client,
	}
}

// Execute creates a new subscription for the given owner/repo.
func (s *Subscribe) Execute(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	return s.client.Subscribe(ctx, owner, repo)
}
