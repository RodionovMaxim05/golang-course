package ratelimiter

import (
	"context"
	"log/slog"
)

type RateLimiter interface {
	Allow(context.Context, string) (bool, float64, error)
}

type FallbackRateLimiter struct {
	log       *slog.Logger
	primary   RateLimiter // Redis RateLimiter
	secondary RateLimiter // In-Memory RateLimiter
}

func NewFallbackRateLimiter(primary, secondary RateLimiter, log *slog.Logger) *FallbackRateLimiter {
	return &FallbackRateLimiter{
		log:       log,
		primary:   primary,
		secondary: secondary,
	}
}

func (frt *FallbackRateLimiter) Allow(ctx context.Context, key string) (bool, float64, error) {
	allowed, remaining, err := frt.primary.Allow(ctx, key)
	if err == nil {
		return allowed, remaining, nil
	}

	frt.log.Debug("primary rate limiter failed. Falling back to in-memory.", "error", err)

	allowedSecondary, remainingSecondary, err := frt.secondary.Allow(ctx, key)
	if err != nil {
		frt.log.Error("secondary rate limiter failed", "error", err)
		return false, 0, err
	}

	return allowedSecondary, remainingSecondary, nil
}
