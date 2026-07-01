package domain

import "errors"

var (
	ErrNotFound    = errors.New("repository not found")
	ErrRateLimited = errors.New("rate limit exceeded")
)
