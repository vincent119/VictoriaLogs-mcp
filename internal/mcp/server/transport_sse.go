package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/vincent119/zlogger"
)

// ServeSSE starts server with SSE transport
func (s *MCPServer) ServeSSE(addr string) error {
	zlogger.Info("MCP Server starting",
		zlogger.String("transport", "sse"),
		zlogger.String("addr", addr),
		zlogger.String("name", s.cfg.Server.Name),
		zlogger.String("version", s.cfg.Server.Version),
	)

	// Create SSE server and start listening
	// NewSSEServer takes (server, ...options), so we don't pass URL as string here.
	sse := server.NewSSEServer(s.server)
	return sse.Start(addr)
}
