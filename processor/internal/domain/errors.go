package domain

import "errors"

// Package-level sentinel errors returned by use cases and storage adapters.
var (
	ErrNotFound        = errors.New("not found")
	ErrAccepted        = errors.New("request accepted")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInternal        = errors.New("internal error")
	ErrRateLimited     = errors.New("rate limit exceeded")
	ErrRepoNotFound    = errors.New("repository not found on GitHub")
)
