package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/dto"
	"repo-watcher/api/internal/usecase"
)

// CreateSubscription godoc
// @Summary Subscribe to the repository
// @Description Subscribe to receive information about the GitHub repository
// @Produce json
// @Param url query string true "GitHub repository URL (e.g. https://github.com/golang/go)"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func NewSubscribeHandler(log *slog.Logger, subscribe *usecase.Subscribe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			writeJSON(w, log, http.StatusBadRequest, map[string]string{"error": "url is required"})
			return
		}

		owner, repo, err := parseGitHubURL(url)
		if err != nil {
			log.Error("failed to parse github url", "error", err)
			writeJSON(w, log, http.StatusBadRequest, map[string]string{"error": "failed to parse github url"})
			return
		}

		subscription, err := subscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to subscribe repository", "error", err)
			writeJSON(w, log, httpCode, map[string]string{"error": err.Error()})
			return
		}

		log.Info("repository subscription created", "owner", owner, "repo", repo)

		response := dto.SubscriptionResponse{
			Owner:     subscription.Owner,
			Repo:      subscription.Repo,
			CreatedAt: subscription.CreatedAt,
		}

		writeJSON(w, log, http.StatusOK, response)
	}
}

// Unsubscribe godoc
// @Summary Unsubscribe from the repository
// @Description Remove subscription for GitHub repository
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{owner}/{repo} [delete]
func NewUnsubscribeHandler(log *slog.Logger, unsubscribe *usecase.Unsubscribe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		if owner == "" || repo == "" {
			writeJSON(w, log, http.StatusBadRequest, map[string]string{"error": "owner and repo are required"})
			return
		}

		err := unsubscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to unsubscribe repository", "error", err)
			writeJSON(w, log, httpCode, map[string]string{"error": err.Error()})
			return
		}

		log.Info("repository unsubscribed", "owner", owner, "repo", repo)

		writeJSON(w, log, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// ListSubscriptions godoc
// @Summary Get current subscription list
// @Description Return all subscribed GitHub repositories
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func NewListSubscriptionsHandler(log *slog.Logger, subscriptions *usecase.GetSubscriptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := subscriptions.Execute(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			httpCode := DomainErrToHTTP(err)
			writeJSON(w, log, httpCode, map[string]string{"error": err.Error()})
			return
		}

		result := make([]dto.SubscriptionResponse, 0, len(resp))
		for _, sub := range resp {
			result = append(result, dto.SubscriptionResponse{
				Owner:     sub.Owner,
				Repo:      sub.Repo,
				CreatedAt: sub.CreatedAt,
			})
		}

		log.Info("subscriptions listed", "count", len(result))

		writeJSON(w, log, http.StatusOK, result)
	}
}

// SubscriptionsInfo godoc
// @Summary Get subscribed repositories info
// @Description Retrieve aggregated information for all subscribed repositories
// @Success 200 {array} dto.RepoInfoResponse
// @Failure 500 {object} map[string]string
// @Router /subscriptions/info [get]
func NewSubscriptionsInfoHandler(log *slog.Logger, subscriptionsInfo *usecase.GetSubscriptionsInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := subscriptionsInfo.Execute(r.Context())
		if err != nil {
			log.Error("failed to fetch subscriptions info", "error", err)
			httpCode := DomainErrToHTTP(err)
			writeJSON(w, log, httpCode, map[string]string{"error": err.Error()})
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
