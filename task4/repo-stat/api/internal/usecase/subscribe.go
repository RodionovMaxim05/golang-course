package usecase

import (
	"context"
)

type Subscriber interface {
	Subscribe(ctx context.Context, owner, repo string) error
}

type Subscribe struct {
	subscriber Subscriber
}

func NewSubscriber(subscriber Subscriber) *Subscribe {
	return &Subscribe{
		subscriber: subscriber,
	}
}

func (s *Subscribe) Execute(ctx context.Context, owner, repo string) error {
	return s.subscriber.Subscribe(ctx, owner, repo)
}
