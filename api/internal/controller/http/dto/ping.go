package dto

// PingResponse is the JSON response body for the health check endpoint.
type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

// ServiceStatus reports the liveness status of a single downstream service.
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
