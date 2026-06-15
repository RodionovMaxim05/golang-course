package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"repo-stat/processor/internal/domain"
)

type DBRepository struct {
	pool    *pgxpool.Pool
	queries Querier
}

func NewDBRepository(pool *pgxpool.Pool) *DBRepository {
	return &DBRepository{
		pool:    pool,
		queries: New(pool),
	}
}

func (r *DBRepository) InsertRepo(ctx context.Context, repo *domain.Repository) error {
	err := r.queries.InsertRepo(ctx, InsertRepoParams{
		FullName:        repo.FullName,
		Description:     pgtype.Text{String: repo.Description, Valid: true},
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       pgtype.Timestamptz{Time: repo.CreatedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("insert repo into db: %w", err)
	}

	return nil
}

func (r *DBRepository) GetRepo(ctx context.Context, fullName string) (*domain.Repository, error) {
	row, err := r.queries.GetRepo(ctx, fullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("get repo from db: %w", err)
	}

	return toRepository(row), nil
}

func (r *DBRepository) GetAllRepos(ctx context.Context) ([]*domain.Repository, error) {
	rows, err := r.queries.ListAllRepos(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all repos from db: %w", err)
	}

	domainRepos := make([]*domain.Repository, 0, len(rows))
	for _, row := range rows {
		domainRepos = append(domainRepos, toRepository(row))
	}
	return domainRepos, nil
}

func toRepository(row Repository) *domain.Repository {
	return &domain.Repository{
		FullName:        row.FullName,
		Description:     row.Description.String,
		StargazersCount: int(row.StargazersCount),
		ForksCount:      int(row.ForksCount),
		CreatedAt:       row.CreatedAt.Time,
	}
}
