package usecase

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type SubscriptionsInfoGetter interface {
	GetSubscriptionsInfo(ctx context.Context) ([]*domain.Repository, error)
}

type GetSubscriptionsInfo struct {
	client SubscriptionsInfoGetter
}

func NewGetSubscriptionsInfo(client SubscriptionsInfoGetter) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{client: client}
}

func (gsi *GetSubscriptionsInfo) Execute(ctx context.Context) ([]*domain.Repository, error) {
	return gsi.client.GetSubscriptionsInfo(ctx)
}
