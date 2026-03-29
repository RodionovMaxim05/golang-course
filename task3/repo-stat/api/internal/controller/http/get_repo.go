package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/api/internal/usecase"
	processorpb "repo-stat/proto/processor"
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
			st := status.Convert(err)
			httpCode := grpcStatusToHTTP(st.Code())

			log.Error("failed to get repo", "error", err, "grpc_code", st.Code())
			writeJSONError(w, httpCode, st.Message())
			return
		}

		response := mapRepoResponse(resp)

		log.Info("repository info fetched successfully", "owner", owner, "repo", repo, "stars", resp.StargazersCount, "forks", resp.ForksCount)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write repo response", "error", err)
		}
	}
}

func grpcStatusToHTTP(code codes.Code) int {
	switch code {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.InvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
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
