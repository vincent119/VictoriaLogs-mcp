// Package logging provides logger initialization
package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init initializes the global logger to stderr (for MCP stdio compatibility)
func Init(level, format string) {
	// Parse level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Create encoder based on format
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Write to stderr for MCP stdio compatibility
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		zapLevel,
	)

	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

// Sync flushes any buffered log entries (call before exit)
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}
