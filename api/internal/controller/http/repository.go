package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/usecase"
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
			writeJSON(w, log, http.StatusBadRequest, map[string]string{"error": "failed to parse github url"})
			return
		}

		repoInfo, err := getRepo.Execute(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to get repo", "error", err)
			httpCode := DomainErrToHTTP(err)
			writeJSON(w, log, httpCode, map[string]string{"error": err.Error()})
			return
		}

		log.Info("repository info fetched successfully", "owner", owner, "repo", repo, "stars", repoInfo.Stargazers, "forks", repoInfo.Forks)

		response := mapRepoResponse(repoInfo)

		writeJSON(w, log, http.StatusOK, response)
	}
}
