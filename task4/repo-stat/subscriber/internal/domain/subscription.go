package domain

import "time"

type Subscription struct {
	Owner string
	Repo  string
}

type SubscriptionResponse struct {
	Owner     string
	Repo      string
	CreatedAt time.Time
}
