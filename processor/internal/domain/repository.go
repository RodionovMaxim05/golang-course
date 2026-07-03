package domain

import (
	"context"
	"time"
)

type Status string

func (s Status) String() string {
	return string(s)
}

func ParseStatus(s string) Status {
	status := Status(s)

	switch status {
	case StatusSuccess, StatusPending, StatusError:
		return status
	default:
		return "UNKNOWN"
	}
}

const (
	StatusSuccess Status = "SUCCESS"
	StatusPending Status = "PENDING"
	StatusError   Status = "ERROR"
)

type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
	Status          Status
	ErrorCode       string
}

type DataStorage interface {
	InsertRepo(ctx context.Context, repo *Repository) error
	UpdateRepoStatus(ctx context.Context, repo *Repository) error
	GetRepo(ctx context.Context, fullName string) (*Repository, error)
	GetReposByNames(ctx context.Context, names []string) ([]*Repository, error)
}
