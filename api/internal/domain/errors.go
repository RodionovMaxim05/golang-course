package domain

import (
	"errors"
)

// Sentinel errors returned by use cases, translated from downstream gRPC
// status codes returned by the Subscriber and Processor services.
var (
	ErrNotFound        = errors.New("not found")
	ErrRateLimited     = errors.New("rate limited")
	ErrUnavailable     = errors.New("service unavailable")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInternal        = errors.New("internal error")
	ErrPending         = errors.New("repository data collection in progress")
	ErrAlreadyExists   = errors.New("already exists")
)
