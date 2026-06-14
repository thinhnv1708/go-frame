package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_MissingAPP_ENV(t *testing.T) {
	t.Setenv("APP_ENV", "")

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "APP_ENV is required")
}

func TestLoad_FileNotFound(t *testing.T) {
	t.Setenv("APP_ENV", "nonexistent")

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config file not found")
}

func TestLoad_Success(t *testing.T) {
	t.Setenv("APP_ENV", "dev")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "identify-dev", cfg.App.Name)
	assert.Equal(t, 8080, cfg.App.Port)
}
