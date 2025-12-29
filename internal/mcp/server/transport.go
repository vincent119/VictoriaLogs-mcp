package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/vincent119/zlogger"
)

// ServeStdio starts server with Stdio transport
func (s *MCPServer) ServeStdio() error {
	zlogger.Info("MCP Server starting",
		zlogger.String("transport", "stdio"),
		zlogger.String("name", s.cfg.Server.Name),
		zlogger.String("version", s.cfg.Server.Version),
	)

	return server.ServeStdio(s.server)
}

// ServeTCP starts server with TCP transport (placeholder)
func (s *MCPServer) ServeTCP(addr string) error {
	zlogger.Info("MCP Server starting",
		zlogger.String("transport", "tcp"),
		zlogger.String("addr", addr),
		zlogger.String("name", s.cfg.Server.Name),
		zlogger.String("version", s.cfg.Server.Version),
	)

	// TODO: Implement TCP transport
	// Currently mcp-go mainly supports stdio, TCP requires additional implementation
	return nil
}
