package domain

import "errors"

var ErrNotFound = errors.New("repository not found")
var ErrRateLimited = errors.New("rate limit exceeded")
