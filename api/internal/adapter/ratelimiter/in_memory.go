package ratelimiter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"repo-watcher/api/config"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type InMemoryRateLimiter struct {
	log        *slog.Logger
	mu         sync.RWMutex
	limiters   map[string]*entry
	rate       rate.Limit
	burst      int
	ttl        time.Duration
	sampleSize int
}

func NewInMemoryRateLimiter(cfg config.RateLimit, log *slog.Logger) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		log:        log,
		limiters:   make(map[string]*entry),
		rate:       rate.Limit(cfg.RequestsPerSecond),
		burst:      cfg.Burst,
		ttl:        5 * time.Minute,
		sampleSize: 20,
	}
}

func (imrl *InMemoryRateLimiter) Allow(_ context.Context, key string) (bool, float64, error) {
	imrl.mu.RLock()
	e, exists := imrl.limiters[key]
	imrl.mu.RUnlock()

	if exists {
		e.lastSeen = time.Now()
		return allowResult(e.limiter)
	}

	imrl.mu.Lock()
	defer imrl.mu.Unlock()

	if e, exists = imrl.limiters[key]; exists {
		e.lastSeen = time.Now()
		return allowResult(e.limiter)
	}

	imrl.limiters[key] = &entry{
		limiter:  rate.NewLimiter(imrl.rate, imrl.burst),
		lastSeen: time.Now(),
	}

	imrl.cleanupSample()

	return allowResult(e.limiter)
}

func allowResult(l *rate.Limiter) (bool, float64, error) {
	allowed := l.Allow()
	remaining := l.Tokens()
	return allowed, remaining, nil
}

func (imrl *InMemoryRateLimiter) cleanupSample() {
	if len(imrl.limiters) == 0 {
		return
	}

	now := time.Now()
	checked := 0

	// Map iteration in Go is randomized
	for key, e := range imrl.limiters {
		if checked >= imrl.sampleSize {
			break
		}
		checked++

		if now.Sub(e.lastSeen) > imrl.ttl {
			delete(imrl.limiters, key)
		}
	}
}
