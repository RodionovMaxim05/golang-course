package domain

import "time"

type Repository struct {
	FullName        string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
}
