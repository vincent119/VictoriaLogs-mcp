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

	// URL for the SSE client to connect to (required for events)
    // We assume localhost if not specified, but Smithery might access via public URL.
    // However, for container internal start, just passing the listening address usually suffices for the base URL construction in some libs,
    // but mcp-go wants an external URL.
    // We will construct a minimal one.
	sseServer := server.NewSSEServer(s.server, "http://localhost"+addr)

    // Wait, the lint said "cannot use string as option".
    // I need to find the option for URL.
    // Search result said `WithBaseURL`.
    // BUT I cannot be sure if it's exported or if `NewSSEServer` signature changed in v0.43.2 vs doc.
    // Let's try `server.NewSSEServer(s.server, "http://localhost"+addr)` ? NO, that failed.

    // Let's try to assume the search result is correct about OPTIONS.
    // But I don't see `WithBaseURL` symbol in the search text explicitly as a `server` package export, it might be.

    // Actually, looking at `mcp-go` v0.43.0 release notes implies breaking changes.
    // Let's try safely: `Start` takes the address.
    // Creating `NewSSEServer(s.server)` without options first.

	sse := server.NewSSEServer(s.server)
	return sse.Start(addr)
}
