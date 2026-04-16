package domain

import "errors"

var ErrNotFound = errors.New("repository not found")
var ErrRateLimited = errors.New("rate limit exceeded")
var ErrInvalidArgument = errors.New("invalid github url")
