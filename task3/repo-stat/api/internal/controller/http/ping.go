package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// Ping godoc
// @Summary Ping services
// @Description Get ping status of processor and subscriber
// @Success 200 {object} dto.PingResponse
// @Failure 503 {object} dto.PingResponse
// @Router /api/ping [get]
func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriberStatus, processorStatus := ping.Execute(r.Context())

		services := []dto.ServiceStatus{
			{Name: "processor", Status: string(processorStatus)},
			{Name: "subscriber", Status: string(subscriberStatus)},
		}

		globalStatus := "ok"
		statusCode := http.StatusOK

		for _, s := range services {
			if s.Status == "down" {
				globalStatus = "degraded"
				statusCode = http.StatusServiceUnavailable
				break
			}
		}

		log.Info("ping request processed", "status", globalStatus, "processor", processorStatus, "subscriber", subscriberStatus)

		response := dto.PingResponse{
			Status:   globalStatus,
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
