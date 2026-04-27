package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type SubscriptionsInfoGetter interface {
	GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error)
}

type GetSubscriptionsInfo struct {
	subsInfoGetter SubscriptionsInfoGetter
}

func NewGetSubscriptionsInfo(subsInfoGetter SubscriptionsInfoGetter) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{
		subsInfoGetter: subsInfoGetter,
	}
}

func (gsi *GetSubscriptionsInfo) Execute(ctx context.Context) ([]domain.Repository, error) {
	return gsi.subsInfoGetter.GetSubscriptionsInfo(ctx)
}
