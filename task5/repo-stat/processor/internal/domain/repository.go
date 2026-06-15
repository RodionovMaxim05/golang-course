package domain

import (
	"context"
	"time"
)

type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
}

type DataStorage interface {
	InsertRepo(ctx context.Context, repo *Repository) error
	GetRepo(ctx context.Context, fullName string) (*Repository, error)
	GetAllRepos(ctx context.Context) ([]*Repository, error)
}
