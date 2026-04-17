package http

import (
	"log/slog"
	"net/http"

	"repo-stat/api/internal/usecase"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	ping *usecase.Ping,
	getRepo *usecase.GetRepo,
	subscribe *usecase.Subscribe,
	unsubscribe *usecase.Unsubscribe,
	getSubscriptions *usecase.GetSubscriptions,
	subscriptionsInfo *usecase.SubscriptionsInfo,
) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewGetRepoHandler(log, getRepo))
	mux.Handle("GET /subscriptions/info", NewSubscriptionsInfoHandler(log, subscriptionsInfo))
	mux.Handle("GET /subscriptions", NewListSubscriptionsHandler(log, getSubscriptions))
	mux.Handle("POST /subscriptions", NewSubscribeHandler(log, subscribe))
	mux.Handle("DELETE /subscriptions/{owner}/{repo}", NewUnsubscribeHandler(log, unsubscribe))
	mux.Handle("/swagger/", SwaggerHandler())
}
