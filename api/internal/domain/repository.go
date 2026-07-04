package domain

import "time"

// Repository is the domain entity representing GitHub repository metrics,
// as returned by the Processor service.
type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
}
