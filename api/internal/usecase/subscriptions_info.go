package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type SubscriptionsInfoGetter interface {
	GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error)
}

type GetSubscriptionsInfo struct {
	client SubscriptionsInfoGetter
}

func NewGetSubscriptionsInfo(client SubscriptionsInfoGetter) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{
		client: client,
	}
}

// Execute returns aggregated repository metrics for all subscriptions.
func (gsi *GetSubscriptionsInfo) Execute(ctx context.Context) ([]domain.Repository, error) {
	return gsi.client.GetSubscriptionsInfo(ctx)
}
