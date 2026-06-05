package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter provides per-key rate limiting.
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

// Limit returns a gin middleware that rate limits by the given key.
func (rl *RateLimiter) Limit(getKey func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getKey(c)
		if key == "" {
			c.Next()
			return
		}

		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		// Filter old requests
		var valid []time.Time
		for _, t := range rl.requests[key] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.limit {
			rl.requests[key] = valid
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		valid = append(valid, now)
		rl.requests[key] = valid
		rl.mu.Unlock()

		c.Next()
	}
}

// cleanup periodically removes old entries.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)
		for key, times := range rl.requests {
			var valid []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}
