// Package app provides application lifecycle management
package app

import (
	"time"

	"github.com/vincent119/victorialogs-mcp/internal/config"
	"github.com/vincent119/victorialogs-mcp/internal/logging"
	mcpserver "github.com/vincent119/victorialogs-mcp/internal/mcp/server"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
	"github.com/vincent119/victorialogs-mcp/internal/util"
	"github.com/vincent119/victorialogs-mcp/internal/victorialogs"
	"github.com/vincent119/zlogger"
)

// Application struct
type Application struct {
	cfg       *config.Config
	mcpServer *mcpserver.MCPServer
	vlClient  *victorialogs.Client
	policyMgr *policy.Manager
}

// New creates a new application
func New(cfg *config.Config) (*Application, error) {
	app := &Application{
		cfg: cfg,
	}

	// Initialize Logger
	logging.Init(cfg.Logging.Level, cfg.Logging.Format)

	// Initialize VictoriaLogs Client
	app.vlClient = victorialogs.NewClient(
		cfg.VictoriaLogs.URL,
		util.AuthConfig{
			Type:     cfg.VictoriaLogs.Auth.Type,
			Username: cfg.VictoriaLogs.Auth.Username,
			Password: cfg.VictoriaLogs.Auth.Password,
			Token:    cfg.VictoriaLogs.Auth.Token,
		},
		cfg.VictoriaLogs.Timeout,
		victorialogs.WithMaxResults(cfg.VictoriaLogs.MaxResults),
	)

	// Initialize Policy Manager
	app.policyMgr = policy.NewManager(policy.Config{
		RateLimit: policy.RateLimitConfig{
			Enabled:           cfg.Policy.RateLimit.Enabled,
			RequestsPerMinute: cfg.Policy.RateLimit.RequestsPerMinute,
		},
		Allowlist: policy.AllowlistConfig{
			Enabled: cfg.Policy.Allowlist.Enabled,
			Streams: cfg.Policy.Allowlist.Streams,
		},
		CircuitBreaker: policy.CircuitBreakerConfig{
			Enabled:        cfg.Policy.CircuitBreaker.Enabled,
			ErrorThreshold: cfg.Policy.CircuitBreaker.ErrorThreshold,
			Timeout:        cfg.Policy.CircuitBreaker.Timeout.String(),
		},
	})

	// Initialize MCP Server
	app.mcpServer = mcpserver.New(cfg, app.vlClient, app.policyMgr)

	zlogger.Info("Application initialized",
		zlogger.String("name", cfg.Server.Name),
		zlogger.String("transport", cfg.Server.Transport),
		zlogger.String("victorialogs_url", cfg.VictoriaLogs.URL),
	)

	return app, nil
}

// Run executes the application
func (app *Application) Run() error {
	switch app.cfg.Server.Transport {
	case "stdio":
		return app.mcpServer.ServeStdio()
	case "tcp":
		return app.mcpServer.ServeTCP(app.cfg.Server.TCPAddr)
	case "sse":
		return app.mcpServer.ServeSSE(app.cfg.Server.TCPAddr)
	default:
		return app.mcpServer.ServeStdio()
	}
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	zlogger.Info("Shutting down application...")

	// Allow time for in-flight requests
	time.Sleep(100 * time.Millisecond)

	// Close MCP Server
	if err := app.mcpServer.Close(); err != nil {
		zlogger.Error("Failed to close MCP Server", zlogger.Err(err))
	}

	// Sync logger
	logging.Sync()

	zlogger.Info("Application shutdown complete")
	return nil
}

// GetConfig returns the configuration
func (app *Application) GetConfig() *config.Config {
	return app.cfg
}
