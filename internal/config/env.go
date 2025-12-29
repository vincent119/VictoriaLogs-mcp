package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load 載入設定檔
func Load(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	v := viper.New()

	// 設定預設值
	setDefaults(v)

	// 設定環境變數
	v.SetEnvPrefix("VLMCP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 載入設定檔
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 嘗試從預設路徑載入
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.AddConfigPath("/etc/vlmcp")
	}

	if err := v.ReadInConfig(); err != nil {
		// 如果找不到設定檔，使用預設值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 解析設定到 struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 驗證設定
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// MustLoad 載入設定，失敗時 panic
func MustLoad(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

// setDefaults 設定 Viper 預設值
func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.name", "victorialogs-mcp")
	v.SetDefault("server.version", "1.0.0")
	v.SetDefault("server.transport", "stdio")
	v.SetDefault("server.tcp_addr", ":9090")

	// VictoriaLogs
	v.SetDefault("victorialogs.url", "http://localhost:9428")
	v.SetDefault("victorialogs.timeout", "30s")
	v.SetDefault("victorialogs.query_timeout", "60s")
	v.SetDefault("victorialogs.max_results", 5000)
	v.SetDefault("victorialogs.auth.type", "none")

	// Policy
	v.SetDefault("policy.rate_limit.enabled", true)
	v.SetDefault("policy.rate_limit.requests_per_minute", 60)
	v.SetDefault("policy.allowlist.enabled", false)
	v.SetDefault("policy.circuit_breaker.enabled", true)
	v.SetDefault("policy.circuit_breaker.error_threshold", 5)
	v.SetDefault("policy.circuit_breaker.timeout", "30s")

	// Logging
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
}

// GetEnv 取得環境變數，支援預設值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
