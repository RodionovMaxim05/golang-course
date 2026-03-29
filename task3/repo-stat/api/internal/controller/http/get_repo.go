package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"repo-stat/api/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			log.Error("failed to get repo", "error", err)
			st := status.Convert(err)
			code := http.StatusInternalServerError
			switch st.Code() {
			case codes.NotFound:
				code = http.StatusNotFound
			case codes.ResourceExhausted:
				code = http.StatusTooManyRequests
			case codes.Unavailable:
				code = http.StatusServiceUnavailable
			case codes.InvalidArgument:
				code = http.StatusBadRequest
			}

			writeJSONError(w, code, st.Message())
			return
		}

		response := mapRepoResponse(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write repo response", "error", err)
		}
	}
}

func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	if u.Host != "github.com" {
		return "", "", fmt.Errorf("invalid host: %s", u.Host)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid path: %s", u.Path)
	}
	return parts[0], parts[1], nil
}

func mapRepoResponse(resp *processorpb.GetRepoResponse) map[string]interface{} {
	return map[string]interface{}{
		"full_name":   resp.Name + "/" + resp.Repo,
		"description": resp.Description,
		"stars":       resp.StargazersCount,
		"forks":       resp.ForksCount,
		"created_at":  resp.CreatedAt.AsTime().Format(time.RFC3339),
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
