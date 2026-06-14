package app

import (
	"context"

	"identify/internal/logger"
	"identify/internal/provider"
)

// App manages the application lifecycle and provider orchestration.
type App struct {
	registry *provider.Registry
	logger   logger.Logger
}

// NewApp creates a new App instance.
func NewApp(registry *provider.Registry, logger logger.Logger) *App {
	return &App{
		registry: registry,
		logger:   logger,
	}
}

// Run boots all providers and blocks until the context is done, then closes all providers.
func (a *App) Run(ctx context.Context) error {
	// Boot all providers
	if err := a.registry.BootAll(ctx); err != nil {
		a.logger.Error("failed to boot providers", "error", err)
		a.registry.CloseAll(context.Background())
		return err
	}

	a.logger.Info("all providers booted successfully")

	// Wait for shutdown signal
	<-ctx.Done()

	a.logger.Info("shutdown signal received")

	// Close all providers (in reverse order of boot)
	if errs := a.registry.CloseAll(context.Background()); len(errs) > 0 {
		for _, err := range errs {
			a.logger.Error("provider close error", "error", err)
		}
	}

	return nil
}
