package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vincent119/zlogger"
)

// WaitForShutdown waits for shutdown signal
func WaitForShutdown(ctx context.Context, app *Application) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		zlogger.Info("Received shutdown signal",
			zlogger.String("signal", sig.String()),
		)
	case <-ctx.Done():
		zlogger.Info("Context cancelled")
	}

	if err := app.Shutdown(); err != nil {
		zlogger.Error("Error during application shutdown",
			zlogger.Err(err),
		)
	}
}

// RunWithGracefulShutdown runs the application with graceful shutdown support
func RunWithGracefulShutdown(app *Application) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start shutdown signal listener
	go WaitForShutdown(ctx, app)

	// Run application
	return app.Run()
}
