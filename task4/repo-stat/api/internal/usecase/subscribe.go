package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Subscriber interface {
	Subscribe(ctx context.Context, owner, repo string) (*domain.Subscription, error)
}

type Subscribe struct {
	subscriber Subscriber
}

func NewSubscriber(subscriber Subscriber) *Subscribe {
	return &Subscribe{
		subscriber: subscriber,
	}
}

func (s *Subscribe) Execute(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	return s.subscriber.Subscribe(ctx, owner, repo)
}
