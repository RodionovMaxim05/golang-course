package domain

import "time"

// Subscription represents a repository being watched by the user, as
// managed by the Subscriber service.
type Subscription struct {
	Owner     string
	Repo      string
	CreatedAt time.Time
}
