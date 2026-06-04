package handler

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// rateCounter tracks one client's request count within the current fixed window.
type rateCounter struct {
	mu          sync.Mutex
	count       int
	windowStart time.Time
}

// RateLimiter enforces a per-IP fixed-window request limit. Counters are kept in
// a sync.Map keyed by client IP; a background sweep evicts idle entries so the
// map does not grow without bound.
type RateLimiter struct {
	clients sync.Map // string(ip) -> *rateCounter
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a limiter allowing limit requests per window per IP and
// starts its background cleanup goroutine.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{limit: limit, window: window}
	go rl.cleanupLoop()
	return rl
}

// allow records a request from ip and reports whether it is within the limit.
func (rl *RateLimiter) allow(ip string) bool {
	now := time.Now()
	v, _ := rl.clients.LoadOrStore(ip, &rateCounter{windowStart: now})
	c := v.(*rateCounter)

	c.mu.Lock()
	defer c.mu.Unlock()
	if now.Sub(c.windowStart) >= rl.window {
		c.count = 0
		c.windowStart = now
	}
	c.count++
	return c.count <= rl.limit
}

// Middleware returns the rate-limiting HTTP middleware.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// cleanupLoop periodically removes counters idle for more than two windows.
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.clients.Range(func(key, value any) bool {
			c := value.(*rateCounter)
			c.mu.Lock()
			stale := now.Sub(c.windowStart) >= 2*rl.window
			c.mu.Unlock()
			if stale {
				rl.clients.Delete(key)
			}
			return true
		})
	}
}

// clientIP extracts the client IP from RemoteAddr (stripping the port).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
