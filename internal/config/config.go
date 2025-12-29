// Package config 提供設定管理功能
package config

import (
	"fmt"
	"time"
)

// Config 應用程式主設定
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	VictoriaLogs  VictoriaLogsConfig  `mapstructure:"victorialogs"`
	Policy        PolicyConfig        `mapstructure:"policy"`
	Logging       LoggingConfig       `mapstructure:"logging"`
}

// ServerConfig MCP Server 設定
type ServerConfig struct {
	Name      string `mapstructure:"name"`
	Version   string `mapstructure:"version"`
	Transport string `mapstructure:"transport"` // stdio | tcp | sse
	TCPAddr   string `mapstructure:"tcp_addr"`
}

// VictoriaLogsConfig VictoriaLogs 連線設定
type VictoriaLogsConfig struct {
	URL          string        `mapstructure:"url"`
	Auth         AuthConfig    `mapstructure:"auth"`
	Timeout      time.Duration `mapstructure:"timeout"`
	QueryTimeout time.Duration `mapstructure:"query_timeout"`
	MaxResults   int           `mapstructure:"max_results"`
}

// AuthConfig 認證設定
type AuthConfig struct {
	Type     string `mapstructure:"type"` // none | basic | bearer
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Token    string `mapstructure:"token"`
}

// PolicyConfig 安全策略設定
type PolicyConfig struct {
	RateLimit      RateLimitConfig      `mapstructure:"rate_limit"`
	Allowlist      AllowlistConfig      `mapstructure:"allowlist"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
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
}

// CircuitBreakerConfig Circuit Breaker 設定
type CircuitBreakerConfig struct {
	Enabled        bool          `mapstructure:"enabled"`
	ErrorThreshold int           `mapstructure:"error_threshold"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

// LoggingConfig 日誌設定
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // debug | info | warn | error
	Format string `mapstructure:"format"` // json | text
}

// Validate 驗證設定
func (c *Config) Validate() error {
	if c.Server.Name == "" {
		return fmt.Errorf("server.name is required")
	}

	if c.Server.Transport != "stdio" && c.Server.Transport != "tcp" && c.Server.Transport != "sse" {
		return fmt.Errorf("server.transport must be 'stdio', 'tcp', or 'sse'")
	}

	if c.Server.Transport == "tcp" && c.Server.TCPAddr == "" {
		return fmt.Errorf("server.tcp_addr is required when transport is 'tcp'")
	}

	if c.VictoriaLogs.URL == "" {
		return fmt.Errorf("victorialogs.url is required")
	}

	if c.VictoriaLogs.Auth.Type != "" &&
		c.VictoriaLogs.Auth.Type != "none" &&
		c.VictoriaLogs.Auth.Type != "basic" &&
		c.VictoriaLogs.Auth.Type != "bearer" {
		return fmt.Errorf("victorialogs.auth.type must be 'none', 'basic', or 'bearer'")
	}

	return nil
}

// DefaultConfig 回傳預設設定
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:      "victorialogs-mcp",
			Version:   "1.0.0",
			Transport: "stdio",
			TCPAddr:   ":9090",
		},
		VictoriaLogs: VictoriaLogsConfig{
			URL:          "http://localhost:9428",
			Timeout:      30 * time.Second,
			QueryTimeout: 60 * time.Second,
			MaxResults:   5000,
			Auth: AuthConfig{
				Type: "none",
			},
		},
		Policy: PolicyConfig{
			RateLimit: RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 60,
			},
			Allowlist: AllowlistConfig{
				Enabled: false,
				Streams: []string{},
			},
			CircuitBreaker: CircuitBreakerConfig{
				Enabled:        true,
				ErrorThreshold: 5,
				Timeout:        30 * time.Second,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
