package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithContext returns a context that is cancelled when SIGINT or SIGTERM is received.
func WithContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	return ctx, cancel
}
