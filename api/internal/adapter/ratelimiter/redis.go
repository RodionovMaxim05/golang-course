package ratelimiter

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log/slog"
	"repo-watcher/api/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// Source: https://github.com/redis/docs/blob/main/content/develop/use-cases/rate-limiter/go/token_bucket.go

// tokenBucketScript is the canonical Lua script for atomic token bucket operations.
// All language implementations use this exact script for behavioral consistency.
const tokenBucketScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local refill_interval = tonumber(ARGV[3])
local now = tonumber(ARGV[4])

-- Get current state or initialize
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

-- Initialize if this is the first request
if tokens == nil then
    tokens = capacity
    last_refill = now
end

-- Calculate token refill
local time_passed = now - last_refill
local refills = math.floor(time_passed / refill_interval)

if refills > 0 then
    tokens = math.min(capacity, tokens + (refills * refill_rate))
    last_refill = last_refill + (refills * refill_interval)
end

-- Try to consume a token
local allowed = 0
if tokens >= 1 then
    tokens = tokens - 1
    allowed = 1
end

-- Update state
redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)

-- Return result: allowed (1 or 0) and remaining tokens
return {allowed, tokens}
`

type RedisRateLimiter struct {
	log            *slog.Logger
	client         *redis.Client
	capacity       int
	refillRate     float64
	refillInterval float64
	scriptSHA      string
	scriptLoaded   bool
}

func NewRedisRateLimiter(client *redis.Client, cfg config.RateLimit, log *slog.Logger) *RedisRateLimiter {
	if cfg.Burst == 0 {
		cfg.Burst = 10
	}
	if cfg.RequestsPerSecond == 0 {
		cfg.RequestsPerSecond = 5
	}

	h := sha1.New()
	h.Write([]byte(tokenBucketScript))
	sha := fmt.Sprintf("%x", h.Sum(nil))

	return &RedisRateLimiter{
		log:            log,
		client:         client,
		capacity:       cfg.Burst,
		refillRate:     float64(cfg.RequestsPerSecond),
		refillInterval: 1.0, // 1 second
		scriptSHA:      sha,
	}
}

// ensureScriptLoaded loads the Lua script into Redis if it hasn't been loaded yet.
func (rrl *RedisRateLimiter) ensureScriptLoaded(ctx context.Context) {
	if !rrl.scriptLoaded {
		sha, err := rrl.client.ScriptLoad(ctx, tokenBucketScript).Result()
		if err != nil {
			rrl.log.Warn("failed to pre-load token bucket lua script into redis", "error", err)
			return
		}

		rrl.scriptSHA = sha
		rrl.scriptLoaded = true
		rrl.log.Info("successfully pre-loaded token bucket lua script", "sha", sha)
	}
}

// Allow checks if a request should be allowed for the given key.
// It atomically checks and updates the token bucket state in Redis.
// Returns whether the request is allowed, the number of remaining tokens,
// and any error encountered.
func (rrl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, float64, error) {
	rrl.ensureScriptLoaded(ctx)

	now := float64(time.Now().UnixMicro()) / 1e6

	args := []interface{}{
		rrl.capacity,
		rrl.refillRate,
		rrl.refillInterval,
		now,
	}

	// Try EVALSHA first (faster if script is cached)
	result, err := rrl.client.EvalSha(ctx, rrl.scriptSHA, []string{key}, args...).Int64Slice()
	if err != nil {
		// Script not in cache, fall back to EVAL
		result, err = rrl.client.Eval(ctx, tokenBucketScript, []string{key}, args...).Int64Slice()
		if err != nil {
			rrl.log.Error("token bucket full eval failed", "key", key, "error", err)
			return false, 0, fmt.Errorf("token bucket eval failed: %w", err)
		}
		rrl.scriptLoaded = false
	}

	allowed := result[0] == 1
	remaining := float64(result[1])

	return allowed, remaining, nil
}
