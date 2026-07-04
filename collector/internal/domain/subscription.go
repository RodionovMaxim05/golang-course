package domain

import "time"

// Subscription represents a single repository being watched by a user,
// as reported by the Subscriber service.
type Subscription struct {
	Owner     string
	Repo      string
	CreatedAt time.Time
}
