package usecase

import (
	"context"
	"repo-watcher/subscriber/internal/domain"
)

type SubscriptionLister interface {
	List(ctx context.Context) ([]domain.SubscriptionRecord, error)
}

type GetSubscriptions struct {
	repository SubscriptionLister
}

func NewGetSubscriptions(repository SubscriptionLister) *GetSubscriptions {
	return &GetSubscriptions{repository: repository}
}

// Execute returns the full list of stored subscriptions.
func (gs *GetSubscriptions) Execute(ctx context.Context) ([]domain.SubscriptionRecord, error) {
	return gs.repository.List(ctx)
}
