package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// GetRepo godoc
// @Summary Get repository info
// @Description Get basic information about a GitHub repository
// @Param url query string true "GitHub repository URL"
// @Success 200 {object} map[string]interface{}
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

		resp, err := getRepo.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to get repo", "error", err)
			writeJSONError(w, httpCode, err.Error())
			return
		}

		response := mapRepoResponse(resp)

		log.Info("repository info fetched successfully", "owner", owner, "repo", repo, "stars", resp.Stargazers, "forks", resp.Forks)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write repo response", "error", err)
		}
	}
}

func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != "github.com" {
		return "", "", fmt.Errorf("unsupported host or invalid url")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid path: %s", u.Path)
	}

	return parts[0], parts[1], nil
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

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
