package domain

import "time"

type Repository struct {
	FullName    string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
}
