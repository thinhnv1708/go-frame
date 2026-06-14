package provider

import "context"

type Provider interface {
	Name() string
	Boot(ctx context.Context) error
	Close(ctx context.Context) error
}
