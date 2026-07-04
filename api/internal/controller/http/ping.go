package http

import (
	"log/slog"
	"net/http"

	"repo-watcher/api/internal/controller/http/dto"
	"repo-watcher/api/internal/domain"
	"repo-watcher/api/internal/usecase"
)

// Ping godoc
// @Summary Ping services
// @Description Get ping status of processor and subscriber
// @Success 200 {object} dto.PingResponse
// @Failure 503 {object} dto.PingResponse
// @Router /api/ping [get]
func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pingResult := ping.Execute(r.Context())

		globalStatus := "ok"
		statusCode := http.StatusOK
		if pingResult.Processor == domain.PingStatusDown || pingResult.Subscriber == domain.PingStatusDown {
			globalStatus = "degraded"
			statusCode = http.StatusServiceUnavailable
		}

		services := []dto.ServiceStatus{
			{Name: "processor", Status: string(pingResult.Processor)},
			{Name: "subscriber", Status: string(pingResult.Subscriber)},
		}

		log.Info("ping request processed", "status", globalStatus, "processor", pingResult.Processor, "subscriber", pingResult.Subscriber)

		writeJSON(w, log, statusCode, dto.PingResponse{Status: globalStatus, Services: services})
	}
}
