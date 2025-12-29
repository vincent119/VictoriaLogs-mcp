package server

import (
	"context"
	"fmt"
	"net"

	"github.com/vincent119/zlogger"
)

// TCPServer TCP transport server (placeholder)
type TCPServer struct {
	addr     string
	listener net.Listener
}

// NewTCPServer creates TCP server
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr: addr,
	}
}

// Start starts TCP server
func (t *TCPServer) Start(ctx context.Context) error {
	// Security check: do not listen on 0.0.0.0
	host, _, err := net.SplitHostPort(t.addr)
	if err != nil {
		return fmt.Errorf("invalid address format: %w", err)
	}

	if host == "" || host == "0.0.0.0" {
		zlogger.Warn("TCP Server should not listen on 0.0.0.0, use 127.0.0.1 or specific IP",
			zlogger.String("addr", t.addr),
		)
	}

	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", t.addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}
	t.listener = listener

	zlogger.Info("TCP Server started",
		zlogger.String("addr", t.addr),
	)

	// TODO: Implement TCP connection handling
	// Currently mcp-go mainly supports stdio transport
	// TCP transport requires additional JSON-RPC over TCP implementation

	<-ctx.Done()
	return ctx.Err()
}

// Stop stops TCP server
func (t *TCPServer) Stop() error {
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}

// GetAddr returns listening address
func (t *TCPServer) GetAddr() string {
	if t.listener != nil {
		return t.listener.Addr().String()
	}
	return t.addr
}

// ValidateTCPAddr validates TCP address security
func ValidateTCPAddr(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address format: %w", err)
	}

	// Disallow listening on all interfaces
	if host == "" || host == "0.0.0.0" {
		return fmt.Errorf("security warning: listening on 0.0.0.0 is not allowed, use 127.0.0.1 or specific private IP")
	}

	// Recommend listening only on localhost
	ip := net.ParseIP(host)
	if ip != nil && !ip.IsLoopback() && !ip.IsPrivate() {
		zlogger.Warn("TCP Server listening on non-private address",
			zlogger.String("addr", addr),
		)
	}

	return nil
}
