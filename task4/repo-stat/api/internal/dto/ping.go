package dto

import "time"

type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SubscriptionResponse struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
}
