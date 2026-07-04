package domain

import "errors"

// Sentinel errors returned by the GitHub client adapter and propagated
// through use cases up to the Kafka result producer, which maps them to
// wire-level error codes.
var (
	ErrNotFound    = errors.New("repository not found")
	ErrRateLimited = errors.New("rate limit exceeded")
)
