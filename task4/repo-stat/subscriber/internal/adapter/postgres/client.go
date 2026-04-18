package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"repo-stat/subscriber/internal/domain"
)

type SubscriptionRepository struct {
	db      *pgxpool.Pool
	queries Querier
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.SubscriptionResponse, error) {
	row, err := r.queries.CreateSubscription(ctx, CreateSubscriptionParams{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	})
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return &domain.SubscriptionResponse{
		Owner:     row.Owner,
		Repo:      row.Repo,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, owner, repo string) error {
	err := r.queries.DeleteSubscription(ctx, DeleteSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]domain.SubscriptionResponse, error) {
	rows, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	subs := make([]domain.SubscriptionResponse, 0, len(rows))
	for _, row := range rows {
		subs = append(subs, domain.SubscriptionResponse{
			Owner:     row.Owner,
			Repo:      row.Repo,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return subs, nil
}

func (r *SubscriptionRepository) Exists(ctx context.Context, owner, repo string) (bool, error) {
	exists, err := r.queries.CheckSubscriptionExists(ctx, CheckSubscriptionExistsParams{Owner: owner, Repo: repo})
	if err != nil {
		return false, fmt.Errorf("check subscription exists: %w", err)
	}
	return exists, nil
}
