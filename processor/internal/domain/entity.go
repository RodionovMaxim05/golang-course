package domain

import (
	"time"
)

// Repository is the domain entity representing cached GitHub repository
// metrics along with its current processing status.
type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
	Status          Status
	ErrorCode       ErrorCode
}
