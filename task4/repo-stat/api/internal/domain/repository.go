package domain

import "time"

type Repository struct {
	Owner       string
	Repo        string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
}
