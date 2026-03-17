package domain

import (
	"time"
)

type Repository struct {
	Name            string
	Description     string
	StargazersCount int
	ForksCount      int
	CreatedAt       time.Time
}
