package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/controller/http/dto"
	"repo-watcher/api/internal/usecase"
)

// NewSubscribeHandler godoc
// @Summary Subscribe to the repository
// @Description Subscribe to receive information about the GitHub repository
// @Produce json
// @Param url query string true "GitHub repository URL (e.g. https://github.com/golang/go)"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse "invalid or missing url"
// @Failure 404 {object} dto.ErrorResponse "repository not found on GitHub"
// @Failure 409 {object} dto.ErrorResponse "already subscribed"
// @Failure 429 {object} dto.ErrorResponse "GitHub API rate limit exceeded"
// @Failure 500 {object} dto.ErrorResponse "internal server error"
// @Router /api/subscriptions [post]
func NewSubscribeHandler(log *slog.Logger, subscribe *usecase.Subscribe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			writeJSON(w, log, http.StatusBadRequest, dto.ErrorResponse{Error: "url is required"})
			return
		}

		owner, repo, err := parseGitHubURL(url)
		if err != nil {
			log.Error("failed to parse github url", "error", err)
			writeJSON(w, log, http.StatusBadRequest, dto.ErrorResponse{Error: "failed to parse github url"})
			return
		}

		subscription, err := subscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to subscribe repository", "error", err)
			writeError(w, log, err)
			return
		}

		log.Info("repository subscription created", "owner", owner, "repo", repo)

		response := mapSubscriptionResponse(*subscription)

		writeJSON(w, log, http.StatusCreated, response)
	}
}

// NewUnsubscribeHandler godoc
// @Summary Unsubscribe from the repository
// @Description Remove subscription for GitHub repository
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse "missing owner or repo"
// @Failure 404 {object} dto.ErrorResponse "subscription not found"
// @Failure 500 {object} dto.ErrorResponse "internal server error"
// @Router /api/subscriptions/{owner}/{repo} [delete]
func NewUnsubscribeHandler(log *slog.Logger, unsubscribe *usecase.Unsubscribe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		if owner == "" || repo == "" {
			writeJSON(w, log, http.StatusBadRequest, dto.ErrorResponse{Error: "owner and repo are required"})
			return
		}

		err := unsubscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to unsubscribe repository", "error", err)
			writeError(w, log, err)
			return
		}

		log.Info("repository unsubscribed", "owner", owner, "repo", repo)

		writeJSON(w, log, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// NewListSubscriptionsHandler godoc
// @Summary Get current subscription list
// @Description Return all subscribed GitHub repositories
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 500 {object} dto.ErrorResponse "internal server error"
// @Router /api/subscriptions [get]
func NewListSubscriptionsHandler(log *slog.Logger, subscriptions *usecase.GetSubscriptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := subscriptions.Execute(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			writeError(w, log, err)
			return
		}

		result := make([]dto.SubscriptionResponse, 0, len(resp))
		for _, sub := range resp {
			result = append(result, mapSubscriptionResponse(sub))
		}

		log.Info("subscriptions listed", "count", len(result))

		writeJSON(w, log, http.StatusOK, result)
	}
}

// NewSubscriptionsInfoHandler godoc
// @Summary Get subscribed repositories info
// @Description Retrieve aggregated information for all subscribed repositories
// @Success 200 {array} dto.RepoInfoResponse
// @Failure 500 {object} dto.ErrorResponse "internal server error"
// @Router /api/subscriptions/info [get]
func NewSubscriptionsInfoHandler(log *slog.Logger, subscriptionsInfo *usecase.GetSubscriptionsInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := subscriptionsInfo.Execute(r.Context())
		if err != nil {
			log.Error("failed to fetch subscriptions info", "error", err)
			writeError(w, log, err)
			return
		}

		result := make([]dto.RepoInfoResponse, 0, len(resp))
		for _, item := range resp {
			result = append(result, mapRepoResponse(item))
		}

		log.Info("subscriptions info fetched", "count", len(result))

		writeJSON(w, log, http.StatusOK, result)
	}
}
