package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"repo-stat/subscriber/internal/domain"
)

type SubscriptionRepository struct {
	pool    *pgxpool.Pool
	queries Querier
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool:    pool,
		queries: New(pool),
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

	resp := toSubscriptionResponse(row.Owner, row.Repo, row.CreatedAt)
	return &resp, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, owner, repo string) error {
	err := r.queries.DeleteSubscription(ctx, DeleteSubscriptionParams{Owner: owner, Repo: repo})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
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
		subs = append(subs, toSubscriptionResponse(row.Owner, row.Repo, row.CreatedAt))
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

func toSubscriptionResponse(owner, repo string, createdAt pgtype.Timestamptz) domain.SubscriptionResponse {
	return domain.SubscriptionResponse{
		Owner:     owner,
		Repo:      repo,
		CreatedAt: createdAt.Time,
	}
}
