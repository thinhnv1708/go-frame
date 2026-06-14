package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithContext_ReturnsContext(t *testing.T) {
	ctx, cancel := WithContext()
	defer cancel()

	// Verify context is not done initially
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done initially")
	default:
		// Expected
	}

	// Cancel via cancel func
	cancel()

	// Wait briefly for cancellation to propagate
	time.Sleep(50 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Expected: context should be cancelled after cancel()
	default:
		t.Fatal("context should be done after cancel()")
	}

	assert.Error(t, ctx.Err())
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestWithContext_CancelFunc(t *testing.T) {
	ctx, cancel := WithContext()
	defer cancel()

	// Verify context is active
	assert.NoError(t, ctx.Err())

	// Cancel
	cancel()
	time.Sleep(50 * time.Millisecond)

	// Verify context is cancelled
	assert.Error(t, ctx.Err())
}
