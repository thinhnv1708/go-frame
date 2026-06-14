package provider

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Registry manages the lifecycle of registered providers.
type Registry struct {
	providers []Provider
	booted    []Provider
	mu        sync.Mutex
}

// NewRegistry creates a new provider Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, p)
}

// BootAll starts all registered providers concurrently using errgroup.
// Tracks successfully-booted providers even when some fail, so CloseAll can clean up.
func (r *Registry) BootAll(ctx context.Context) error {
	booted := make([]bool, len(r.providers))

	g, ctx := errgroup.WithContext(ctx)

	for i, p := range r.providers {
		i, p := i, p
		g.Go(func() error {
			if err := p.Boot(ctx); err != nil {
				return err
			}
			r.mu.Lock()
			booted[i] = true
			r.mu.Unlock()
			return nil
		})
	}

	err := g.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()
	// Populate r.booted with providers that successfully booted
	for i, b := range booted {
		if b {
			r.booted = append(r.booted, r.providers[i])
		}
	}

	return err
}

// CloseAll stops all booted providers in reverse order.
// Close errors are collected and returned, but do not stop the shutdown process.
func (r *Registry) CloseAll(ctx context.Context) []error {
	r.mu.Lock()
	bootedList := make([]Provider, len(r.booted))
	copy(bootedList, r.booted)
	r.mu.Unlock()

	var errs []error
	for i := len(bootedList) - 1; i >= 0; i-- {
		if err := bootedList[i].Close(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
