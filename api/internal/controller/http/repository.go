package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/controller/http/dto"
	"repo-watcher/api/internal/usecase"
)

// GetRepo godoc
// @Summary Get repository info
// @Description Get basic information about a GitHub repository. If the
// @Description repository's data has not been collected yet, returns 202
// @Description with an error message indicating the data is still being
// @Description processed; retry the request later.
// @Param url query string true "GitHub repository URL (e.g. https://github.com/golang/go)"
// @Success 200 {object} dto.RepoInfoResponse
// @Success 202 {object} map[string]string "repository data collection in progress, retry later"
// @Failure 400 {object} map[string]string "invalid or missing url"
// @Failure 404 {object} map[string]string "repository not found on GitHub"
// @Failure 429 {object} map[string]string "GitHub API rate limit exceeded"
// @Failure 500 {object} map[string]string "internal server error"
// @Failure 503 {object} map[string]string "processor service unavailable"
// @Router /api/repositories/info [get]
func NewGetRepoHandler(log *slog.Logger, getRepo *usecase.GetRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner, repo, err := parseGitHubURL(r.URL.Query().Get("url"))
		if err != nil {
			log.Error("failed to parse github url", "error", err)
			writeJSON(w, log, http.StatusBadRequest, dto.ErrorResponse{Error: "failed to parse github url"})
			return
		}

		repoInfo, err := getRepo.Execute(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to get repo", "error", err)
			writeError(w, log, err)
			return
		}

		log.Info(
			"repository info fetched successfully",
			"owner", owner,
			"repo", repo,
			"stars", repoInfo.StargazersCount,
			"forks", repoInfo.ForksCount,
		)

		response := mapRepoResponse(repoInfo)

		writeJSON(w, log, http.StatusOK, response)
	}
}
