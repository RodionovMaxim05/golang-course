package http

import (
	"repo-watcher/api/internal/controller/http/dto"
	"repo-watcher/api/internal/domain"
)

// mapRepoResponse converts a domain.Repository into its JSON response
// representation.
func mapRepoResponse(resp domain.Repository) dto.RepoInfoResponse {
	return dto.RepoInfoResponse{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.StargazersCount,
		Forks:       resp.ForksCount,
		CreatedAt:   resp.CreatedAt,
	}
}

// mapSubscriptionResponse converts a domain.Subscription into its JSON
// response representation.
func mapSubscriptionResponse(sub domain.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		Owner:     sub.Owner,
		Repo:      sub.Repo,
		CreatedAt: sub.CreatedAt,
	}
}
