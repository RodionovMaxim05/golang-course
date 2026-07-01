package usecase

import (
	"context"

	"repo-watcher/subscriber/internal/domain"
)

type GetSubscriptions struct {
	repository domain.SubscriptionRepository
}

func NewGetSubscriptions(repository domain.SubscriptionRepository) *GetSubscriptions {
	return &GetSubscriptions{repository: repository}
}

func (gs *GetSubscriptions) Execute(ctx context.Context) ([]domain.SubscriptionResponse, error) {
	return gs.repository.List(ctx)
}
