package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/controller/http/middleware"
	"repo-watcher/api/internal/usecase"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	rateLimiter middleware.RateLimiter,
	ping *usecase.Ping,
	getRepo *usecase.GetRepo,
	subscribe *usecase.Subscribe,
	unsubscribe *usecase.Unsubscribe,
	getSubscriptions *usecase.GetSubscriptions,
	subscriptionsInfo *usecase.GetSubscriptionsInfo,
) http.Handler {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewGetRepoHandler(log, getRepo))
	mux.Handle("GET /api/subscriptions/info", NewSubscriptionsInfoHandler(log, subscriptionsInfo))
	mux.Handle("GET /api/subscriptions", NewListSubscriptionsHandler(log, getSubscriptions))
	mux.Handle("POST /api/subscriptions", NewSubscribeHandler(log, subscribe))
	mux.Handle("DELETE /api/subscriptions/{owner}/{repo}", NewUnsubscribeHandler(log, unsubscribe))
	mux.Handle("/swagger/", SwaggerHandler())

	return middleware.RateLimit(rateLimiter, log)(mux)
}
