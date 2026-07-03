package http

import (
	"repo-watcher/api/internal/domain"
	"repo-watcher/api/internal/dto"
	"time"
)

func mapRepoResponse(resp domain.Repository) dto.RepoInfoResponse {
	return dto.RepoInfoResponse{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.Stargazers,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt.Format(time.RFC3339),
	}
}
