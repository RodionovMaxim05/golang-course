package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/controller/http/middleware"
	"repo-watcher/api/internal/usecase"
)

// NewHandler builds the complete HTTP handler for the API gateway,
// wiring all routes and the rate-limiting middleware.
func NewHandler(
	log *slog.Logger,
	rateLimiter middleware.RateLimiter,
	pingUC *usecase.Ping,
	getRepoUC *usecase.GetRepo,
	subscribeUC *usecase.Subscribe,
	unsubscribeUC *usecase.Unsubscribe,
	getSubscriptionsUC *usecase.GetSubscriptions,
	subscriptionsInfoUC *usecase.GetSubscriptionsInfo,
) http.Handler {
	mux := http.NewServeMux()
	return AddRoutes(
		mux,
		log,
		rateLimiter,
		pingUC,
		getRepoUC,
		subscribeUC,
		unsubscribeUC,
		getSubscriptionsUC,
		subscriptionsInfoUC,
	)
}
