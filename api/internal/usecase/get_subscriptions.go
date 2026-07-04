package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type SubscriptionLister interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type GetSubscriptions struct {
	client SubscriptionLister
}

func NewGetSubscriptions(client SubscriptionLister) *GetSubscriptions {
	return &GetSubscriptions{
		client: client,
	}
}

// Execute returns all currently active subscriptions.
func (gs *GetSubscriptions) Execute(ctx context.Context) ([]domain.Subscription, error) {
	return gs.client.GetSubscriptions(ctx)
}
