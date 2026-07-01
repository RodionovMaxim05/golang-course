package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"repo-watcher/processor/internal/domain"
)

type RepositoryAdapter struct {
	pool    *pgxpool.Pool
	queries Querier
}

func NewDBRepository(pool *pgxpool.Pool) *RepositoryAdapter {
	return &RepositoryAdapter{
		pool:    pool,
		queries: New(pool),
	}
}

func (r *RepositoryAdapter) InsertRepo(ctx context.Context, repo *domain.Repository) error {
	err := r.queries.InsertRepo(ctx, InsertRepoParams{
		FullName:        repo.FullName,
		Description:     pgtype.Text{String: repo.Description, Valid: true},
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       pgtype.Timestamptz{Time: repo.CreatedAt, Valid: true},
		RepoStatus:      repo.Status,
		ErrorCode:       pgtype.Text{String: repo.ErrorCode, Valid: repo.ErrorCode != ""},
	})
	if err != nil {
		return fmt.Errorf("insert/upsert repo into db: %w", err)
	}

	return nil
}

func (r *RepositoryAdapter) UpdateRepoStatus(ctx context.Context, repo *domain.Repository) error {
	err := r.queries.UpdateRepoStatus(ctx, UpdateRepoStatusParams{
		FullName:   repo.FullName,
		RepoStatus: repo.Status,
		ErrorCode:  pgtype.Text{String: repo.ErrorCode, Valid: repo.ErrorCode != ""},
	})
	if err != nil {
		return fmt.Errorf("update repo status in db: %w", err)
	}

	return nil
}

func (r *RepositoryAdapter) GetRepo(ctx context.Context, fullName string) (*domain.Repository, error) {
	row, err := r.queries.GetRepo(ctx, fullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("get repo from db: %w", err)
	}

	return toRepository(row), nil
}

func (r *RepositoryAdapter) GetReposByNames(ctx context.Context, fullNames []string) ([]*domain.Repository, error) {
	if len(fullNames) == 0 {
		return []*domain.Repository{}, nil
	}

	rows, err := r.queries.GetReposByNames(ctx, fullNames)
	if err != nil {
		return nil, fmt.Errorf("db get repos by names: %w", err)
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
		Status:          row.RepoStatus,
		ErrorCode:       row.ErrorCode.String,
	}
}
