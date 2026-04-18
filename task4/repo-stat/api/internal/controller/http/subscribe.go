package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

type subscribeRequest struct {
	URL string `json:"url"`
}

// CreateSubscription godoc
// @Summary Subscribe to the repository
// @Description Subscribe to receive information about the GitHub repository
// @Accept json
// @Produce json
// @Param request body subscribeRequest true "GitHub repository URL"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func NewSubscribeHandler(log *slog.Logger, subscribe *usecase.Subscribe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req subscribeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode subscribe request", "error", err)
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		owner, repo, err := parseGitHubURL(req.URL)
		if err != nil {
			log.Error("failed to parse github url", "error", err)
			writeJSONError(w, http.StatusBadRequest, "failed to parse github url")
			return
		}

		subscription, err := subscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to subscribe repository", "error", err)
			writeJSONError(w, httpCode, err.Error())
			return
		}

		log.Info("repository subscription created", "owner", owner, "repo", repo)

		response := dto.SubscriptionResponse{
			Owner:     subscription.Owner,
			Repo:      subscription.Repo,
			CreatedAt: subscription.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
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
			writeJSONError(w, http.StatusBadRequest, "owner and repo are required")
			return
		}

		err := unsubscribe.Execute(r.Context(), owner, repo)
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to unsubscribe repository", "error", err)
			writeJSONError(w, httpCode, err.Error())
			return
		}

		log.Info("repository unsubscribed", "owner", owner, "repo", repo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to list subscriptions", "error", err)
			writeJSONError(w, httpCode, err.Error())
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("failed to write subscriptions response", "error", err)
		}
	}
}

// SubscriptionsInfo godoc
// @Summary Get subscribed repositories info
// @Description Retrieve aggregated information for all subscribed repositories
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /subscriptions/info [get]
func NewSubscriptionsInfoHandler(log *slog.Logger, subscriptionsInfo *usecase.SubscriptionsInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := subscriptionsInfo.Execute(r.Context())
		if err != nil {
			httpCode := DomainErrToHTTP(err)
			log.Error("failed to fetch subscriptions info", "error", err)
			writeJSONError(w, httpCode, err.Error())
			return
		}

		result := make([]map[string]interface{}, 0, len(resp))
		for _, item := range resp {
			result = append(result, mapRepoResponse(item))
		}

		log.Info("subscriptions info fetched", "count", len(result))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("failed to write subscriptions info response", "error", err)
		}
	}
}
