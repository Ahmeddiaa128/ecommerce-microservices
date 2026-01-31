package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

// RateLimiter implements a simple rate limiting middleware
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	requests int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		requests: requests,
		window:   window,
	}

	// Clean up old visitors periodically
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{lastSeen: time.Now(), count: 0}
		rl.visitors[ip] = v
	}

	return v
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		v := rl.getVisitor(ip)

		rl.mu.Lock()
		// Reset counter if window has passed
		if time.Since(v.lastSeen) > rl.window {
			v.count = 0
			v.lastSeen = time.Now()
		}

		// Check if limit exceeded
		if v.count >= rl.requests {
			rl.mu.Unlock()
			writeJSONError(c, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		v.count++
		rl.mu.Unlock()

		c.Next()
	}
}
