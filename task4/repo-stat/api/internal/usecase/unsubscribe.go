package usecase

import (
	"context"
)

type Unsubscriber interface {
	Unsubscribe(ctx context.Context, owner, repo string) error
}

type Unsubscribe struct {
	subscriber Unsubscriber
}

func NewUnsubscriber(subscriber Unsubscriber) *Unsubscribe {
	return &Unsubscribe{
		subscriber: subscriber,
	}
}

func (u *Unsubscribe) Execute(ctx context.Context, owner, repo string) error {
	return u.subscriber.Unsubscribe(ctx, owner, repo)
}
