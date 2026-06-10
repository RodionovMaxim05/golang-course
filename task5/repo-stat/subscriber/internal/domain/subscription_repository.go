package domain

import "context"

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) (*SubscriptionResponse, error)
	Delete(ctx context.Context, owner, repo string) error
	List(ctx context.Context) ([]SubscriptionResponse, error)
	Exists(ctx context.Context, owner, repo string) (bool, error)
}
