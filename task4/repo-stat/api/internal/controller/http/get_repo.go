package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// GetRepo godoc
// @Summary Get repository info
// @Description Get basic information about a GitHub repository
// @Param url query string true "GitHub repository URL (e.g. https://github.com/golang/go)"
// @Success 200 {object} dto.RepoInfoResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/repositories/info [get]
func NewGetRepoHandler(log *slog.Logger, getRepo *usecase.GetRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner, repo, err := parseGitHubURL(r.URL.Query().Get("url"))
		if err != nil {
			log.Error("failed to parse github url", "error", err)
			writeJSONError(w, http.StatusBadRequest, "failed to parse github url")
			return
		}

		repoInfo, err := getRepo.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to get repo", "error", err)
			writeJSONError(w, httpCode, err.Error())
			return
		}

		log.Info("repository info fetched successfully", "owner", owner, "repo", repo, "stars", repoInfo.Stargazers, "forks", repoInfo.Forks)

		response := mapRepoResponse(repoInfo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write repo response", "error", err)
		}
	}
}

func mapRepoResponse(resp domain.Repository) dto.RepoInfoResponse {
	return dto.RepoInfoResponse{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.Stargazers,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt.Format(time.RFC3339),
	}
}
