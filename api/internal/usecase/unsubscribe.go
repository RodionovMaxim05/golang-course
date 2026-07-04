package usecase

import (
	"context"
)

type SubscriptionDeleter interface {
	Unsubscribe(ctx context.Context, owner, repo string) error
}

type Unsubscribe struct {
	client SubscriptionDeleter
}

func NewUnsubscriber(client SubscriptionDeleter) *Unsubscribe {
	return &Unsubscribe{
		client: client,
	}
}

// Execute removes the subscription for the given owner/repo.
func (u *Unsubscribe) Execute(ctx context.Context, owner, repo string) error {
	return u.client.Unsubscribe(ctx, owner, repo)
}
