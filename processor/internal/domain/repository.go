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
	Status          string
	ErrorCode       string
}

type DataStorage interface {
	InsertRepo(ctx context.Context, repo *Repository) error
	UpdateRepoStatus(ctx context.Context, repo *Repository) error
	GetRepo(ctx context.Context, fullName string) (*Repository, error)
	GetReposByNames(ctx context.Context, names []string) ([]*Repository, error)
}
