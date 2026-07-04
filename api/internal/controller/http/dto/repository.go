package dto

import "time"

// RepoInfoResponse is the JSON response body for repository information
// endpoints.
type RepoInfoResponse struct {
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
}
