package http

import (
	"log/slog"
	"net/http"

	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, getRepo *usecase.GetRepo) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewGetRepoHandler(log, getRepo))
	mux.Handle("/swagger/", SwaggerHandler())
}
