package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"identify/internal/logger"
	"identify/internal/provider"
)

type mockProvider struct {
	name    string
	bootErr error
	closed  bool
}

func (m *mockProvider) Name() string                   { return m.name }
func (m *mockProvider) Boot(ctx context.Context) error { return m.bootErr }
func (m *mockProvider) Close(ctx context.Context) error {
	m.closed = true
	return nil
}

func TestApp_Run_BootSuccess(t *testing.T) {
	p := &mockProvider{name: "test"}
	r := provider.NewRegistry()
	r.Register(p)

	logger, _ := logger.NewZapLogger("debug", "console")
	app := NewApp(r, logger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- app.Run(ctx)
	}()

	// Give app time to boot
	time.Sleep(100 * time.Millisecond)

	// Cancel to trigger shutdown
	cancel()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}

	assert.True(t, p.closed)
}

func TestApp_Run_BootFailure(t *testing.T) {
	p1 := &mockProvider{name: "p1"}                                      // succeeds
	p2 := &mockProvider{name: "p2", bootErr: errors.New("boot failed")} // fails
	r := provider.NewRegistry()
	r.Register(p1)
	r.Register(p2)

	logger, _ := logger.NewZapLogger("debug", "console")
	app := NewApp(r, logger)

	ctx := context.Background()
	err := app.Run(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "boot failed")
	assert.True(t, p1.closed) // p1 was booted successfully, should be closed
}
