package usecase

import (
	"context"
)

type SubscriptionDeleter interface {
	Delete(ctx context.Context, owner, repo string) error
}

type Unsubscribe struct {
	repository SubscriptionDeleter
}

func NewUnsubscribe(repo SubscriptionDeleter) *Unsubscribe {
	return &Unsubscribe{repository: repo}
}

// Execute removes the subscription for the given owner/repo, if it exists.
func (u *Unsubscribe) Execute(ctx context.Context, owner, repo string) error {
	return u.repository.Delete(ctx, owner, repo)
}
