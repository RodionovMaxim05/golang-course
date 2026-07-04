package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"repo-watcher/subscriber/internal/domain"
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

// Create persists a new subscription. Returns domain.ErrAlreadySubscribed
// if a subscription for the same owner/repo already exists (enforced by
// a unique constraint at the database level).
func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.SubscriptionRecord, error) {
	row, err := r.queries.CreateSubscription(ctx, CreateSubscriptionParams{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrAlreadySubscribed
		}
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	record := toSubscriptionRecord(row.Owner, row.Repo, row.CreatedAt)
	return &record, nil
}

// Delete removes the subscription for the given owner/repo. Returns
// domain.ErrNotFound if no matching subscription exists.
func (r *SubscriptionRepository) Delete(ctx context.Context, owner, repo string) error {
	rowsAffected, err := r.queries.DeleteSubscription(ctx, DeleteSubscriptionParams{Owner: owner, Repo: repo})
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// List returns all currently stored subscriptions.
func (r *SubscriptionRepository) List(ctx context.Context) ([]domain.SubscriptionRecord, error) {
	rows, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	subs := make([]domain.SubscriptionRecord, 0, len(rows))
	for _, row := range rows {
		subs = append(subs, toSubscriptionRecord(row.Owner, row.Repo, row.CreatedAt))
	}
	return subs, nil
}

func toSubscriptionRecord(owner, repo string, createdAt pgtype.Timestamptz) domain.SubscriptionRecord {
	return domain.SubscriptionRecord{
		Owner:     owner,
		Repo:      repo,
		CreatedAt: createdAt.Time,
	}
}
