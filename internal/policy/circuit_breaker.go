package policy

import (
	"fmt"
	"sync"
	"time"
)

// CircuitBreaker implementation
type CircuitBreaker struct {
	enabled        bool
	errorThreshold int
	timeout        time.Duration

	state          CircuitState
	failureCount   int
	lastFailure    time.Time
	halfOpenTries  int
	mu             sync.RWMutex
}

// CircuitState circuit breaker state
type CircuitState int

const (
	// StateClosed normal state
	StateClosed CircuitState = iota
	// StateOpen circuit open state
	StateOpen
	// StateHalfOpen half-open state
	StateHalfOpen
)

// ErrCircuitOpen circuit breaker is open
var ErrCircuitOpen = fmt.Errorf("service temporarily unavailable, please try again later")

// NewCircuitBreaker creates circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	timeout, _ := time.ParseDuration(cfg.Timeout)
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &CircuitBreaker{
		enabled:        cfg.Enabled,
		errorThreshold: cfg.ErrorThreshold,
		timeout:        timeout,
		state:          StateClosed,
	}
}

// Allow checks if request is allowed
func (cb *CircuitBreaker) Allow() error {
	if !cb.enabled {
		return nil
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if can transition to half-open state
		if time.Since(cb.lastFailure) >= cb.timeout {
			cb.state = StateHalfOpen
			cb.halfOpenTries = 0
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Half-open state only allows limited requests
		if cb.halfOpenTries >= 1 {
			return ErrCircuitOpen
		}
		cb.halfOpenTries++
		return nil
	}

	return nil
}

// RecordSuccess records success
func (cb *CircuitBreaker) RecordSuccess() {
	if !cb.enabled {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHalfOpen:
		// Half-open success, close circuit breaker
		cb.state = StateClosed
		cb.failureCount = 0
	case StateClosed:
		// Closed state success, reset failure count
		cb.failureCount = 0
	}
}

// RecordFailure records failure
func (cb *CircuitBreaker) RecordFailure() {
	if !cb.enabled {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failureCount >= cb.errorThreshold {
			cb.state = StateOpen
		}
	case StateHalfOpen:
		// Half-open failure, reopen circuit breaker
		cb.state = StateOpen
	}
}

// GetState gets current state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStateString gets state string
func (cb *CircuitBreaker) GetStateString() string {
	switch cb.GetState() {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Reset resets circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failureCount = 0
	cb.halfOpenTries = 0
}
