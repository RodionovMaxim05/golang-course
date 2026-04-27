package usecase

import (
	"context"

	"repo-stat/subscriber/internal/domain"
)

type Unsubscribe struct {
	repository domain.SubscriptionRepository
}

func NewUnsubscribe(repo domain.SubscriptionRepository) *Unsubscribe {
	return &Unsubscribe{repository: repo}
}

func (u *Unsubscribe) Execute(ctx context.Context, owner, repo string) error {
	return u.repository.Delete(ctx, owner, repo)
}
