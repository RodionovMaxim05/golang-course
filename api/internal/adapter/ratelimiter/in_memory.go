package ratelimiter

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

	"repo-watcher/api/config"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64
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

// Allow checks whether a request for the given key is allowed under its
// token bucket limiter, creating a new limiter on first use.
func (imrl *InMemoryRateLimiter) Allow(_ context.Context, key string) (bool, float64, error) {
	imrl.mu.RLock()
	e, exists := imrl.limiters[key]
	imrl.mu.RUnlock()

	if exists {
		e.lastSeen.Store(time.Now().UnixNano())
		return allowResult(e.limiter)
	}

	imrl.mu.Lock()
	defer imrl.mu.Unlock()

	if e, exists = imrl.limiters[key]; exists {
		e.lastSeen.Store(time.Now().UnixNano())
		return allowResult(e.limiter)
	}

	newEntry := &entry{limiter: rate.NewLimiter(imrl.rate, imrl.burst)}
	newEntry.lastSeen.Store(time.Now().UnixNano())
	imrl.limiters[key] = newEntry

	imrl.cleanupSample()

	return allowResult(newEntry.limiter)
}

func allowResult(l *rate.Limiter) (bool, float64, error) {
	allowed := l.Allow()
	remaining := l.Tokens()
	return allowed, remaining, nil
}

// cleanupSample opportunistically evicts a bounded sample of stale
// entries (unused for longer than the configured TTL) to bound memory
// growth without scanning the entire map on every request.
func (imrl *InMemoryRateLimiter) cleanupSample() {
	if len(imrl.limiters) == 0 {
		return
	}

	checked := 0

	// Map iteration in Go is randomized
	for key, e := range imrl.limiters {
		if checked >= imrl.sampleSize {
			break
		}
		checked++

		if time.Since(time.Unix(0, e.lastSeen.Load())) > imrl.ttl {
			delete(imrl.limiters, key)
		}
	}
}
