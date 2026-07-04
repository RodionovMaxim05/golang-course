package domain

import "time"

// Subscription is the input for creating a new subscription.
type Subscription struct {
	Owner string
	Repo  string
}

// SubscriptionRecord is a persisted subscription with its creation timestamp.
type SubscriptionRecord struct {
	Owner     string
	Repo      string
	CreatedAt time.Time
}
