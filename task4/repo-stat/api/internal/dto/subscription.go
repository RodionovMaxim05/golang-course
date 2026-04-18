package dto

import "time"

type SubscriptionResponse struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
}
