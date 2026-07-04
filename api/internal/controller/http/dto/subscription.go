package dto

import "time"

// SubscriptionResponse is the JSON response body representing a single
// subscription.
type SubscriptionResponse struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
}
