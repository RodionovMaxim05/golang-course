package domain

import "time"

// Repository is the domain entity representing repository metrics fetched
// directly from the GitHub API.
type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
}
