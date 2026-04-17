package domain

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrRateLimited     = errors.New("rate limited")
	ErrUnavailable     = errors.New("service unavailable")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInternal        = errors.New("internal error")
)
