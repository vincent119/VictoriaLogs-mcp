// Package policy 提供安全策略管理
package policy

import (
	"context"
)

// Manager 策略管理器
type Manager struct {
	allowlist      *Allowlist
	rateLimit      *RateLimiter
	redact         *Redactor
	circuitBreaker *CircuitBreaker

}

// Config 策略設定
type Config struct {
	RateLimit      RateLimitConfig      `mapstructure:"rate_limit"`
	Allowlist      AllowlistConfig      `mapstructure:"allowlist"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Redact         RedactConfig         `mapstructure:"redact"`
}

// RateLimitConfig Rate Limit 設定
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
}

// AllowlistConfig Allowlist 設定
type AllowlistConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Streams []string `mapstructure:"streams"`
	Deny    []string `mapstructure:"deny"`
}

// CircuitBreakerConfig Circuit Breaker 設定
type CircuitBreakerConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	ErrorThreshold int    `mapstructure:"error_threshold"`
	Timeout        string `mapstructure:"timeout"`
}

// RedactConfig Redact 設定
type RedactConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Patterns []RedactPattern `mapstructure:"patterns"`
}

// RedactPattern Redact 規則
type RedactPattern struct {
	Name        string `mapstructure:"name"`
	Pattern     string `mapstructure:"pattern"`
	Replacement string `mapstructure:"replacement"`
}

// NewManager 建立策略管理器
func NewManager(cfg Config) *Manager {
	m := &Manager{}

	if cfg.Allowlist.Enabled {
		m.allowlist = NewAllowlist(cfg.Allowlist)
	}

	if cfg.RateLimit.Enabled {
		m.rateLimit = NewRateLimiter(cfg.RateLimit)
	}

	if cfg.Redact.Enabled {
		m.redact = NewRedactor(cfg.Redact)
	}

	if cfg.CircuitBreaker.Enabled {
		m.circuitBreaker = NewCircuitBreaker(cfg.CircuitBreaker)
	}

	return m
}

// CheckAllowlist 檢查 Allowlist
func (m *Manager) CheckAllowlist(_ context.Context, stream string) error {
	if m.allowlist == nil {
		return nil
	}
	return m.allowlist.Check(stream)
}

// CheckRateLimit 檢查 Rate Limit
func (m *Manager) CheckRateLimit(_ context.Context, key string) error {
	if m.rateLimit == nil {
		return nil
	}
	return m.rateLimit.Allow(key)
}

// Redact 執行敏感資訊遮罩
func (m *Manager) Redact(data string) string {
	if m.redact == nil {
		return data
	}
	return m.redact.Apply(data)
}

// CheckCircuitBreaker 檢查 Circuit Breaker
func (m *Manager) CheckCircuitBreaker(_ context.Context) error {
	if m.circuitBreaker == nil {
		return nil
	}
	return m.circuitBreaker.Allow()
}

// RecordSuccess 記錄成功
func (m *Manager) RecordSuccess() {
	if m.circuitBreaker != nil {
		m.circuitBreaker.RecordSuccess()
	}
}

// RecordFailure 記錄失敗
func (m *Manager) RecordFailure() {
	if m.circuitBreaker != nil {
		m.circuitBreaker.RecordFailure()
	}
}

// Close 關閉管理器
func (m *Manager) Close() {
	// 清理資源
}
