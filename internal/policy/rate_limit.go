package policy

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter implementation
type RateLimiter struct {
	enabled           bool
	requestsPerMinute int
	counters          map[string]*rateLimitCounter
	mu                sync.RWMutex
}

type rateLimitCounter struct {
	count     int
	resetTime time.Time
}

// ErrRateLimitExceeded rate limit exceeded
var ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded, please try again later")

// NewRateLimiter creates rate limiter
func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		enabled:           cfg.Enabled,
		requestsPerMinute: cfg.RequestsPerMinute,
		counters:          make(map[string]*rateLimitCounter),
	}
}

// Allow checks if request is allowed
func (r *RateLimiter) Allow(key string) error {
	if !r.enabled {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	counter, exists := r.counters[key]
	if !exists {
		r.counters[key] = &rateLimitCounter{
			count:     1,
			resetTime: now.Add(time.Minute),
		}
		return nil
	}

	// Reset counter if past reset time
	if now.After(counter.resetTime) {
		counter.count = 1
		counter.resetTime = now.Add(time.Minute)
		return nil
	}

	// Check if limit exceeded
	if counter.count >= r.requestsPerMinute {
		return ErrRateLimitExceeded
	}

	counter.count++
	return nil
}

// GetRemaining gets remaining request count
func (r *RateLimiter) GetRemaining(key string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counter, exists := r.counters[key]
	if !exists {
		return r.requestsPerMinute
	}

	if time.Now().After(counter.resetTime) {
		return r.requestsPerMinute
	}

	remaining := r.requestsPerMinute - counter.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset resets counter
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.counters, key)
}

// Cleanup cleans up expired counters
func (r *RateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for key, counter := range r.counters {
		if now.After(counter.resetTime) {
			delete(r.counters, key)
		}
	}
}
