package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type SubscriptionsGetter interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type GetSubscriptions struct {
	subscriberClient SubscriptionsGetter
}

func NewGetSubscriptions(subscriber SubscriptionsGetter) *GetSubscriptions {
	return &GetSubscriptions{
		subscriberClient: subscriber,
	}
}

func (gS *GetSubscriptions) Execute(ctx context.Context) ([]domain.Subscription, error) {
	return gS.subscriberClient.GetSubscriptions(ctx)
}
