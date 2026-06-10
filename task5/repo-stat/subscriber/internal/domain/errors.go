package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInternal          = errors.New("internal error")
	ErrRateLimited       = errors.New("rate limit exceeded")
	ErrRepoNotFound      = errors.New("repository not found on GitHub")
)
