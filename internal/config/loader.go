package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Load reads APP_ENV, loads the corresponding YAML config file and .env file,
// then applies environment overrides and returns the Config struct.
func Load() (*Config, error) {
	// Load environment variables from .env if it exists
	_ = gotenv.Load()

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return nil, fmt.Errorf("APP_ENV is required. Set it to 'dev', 'prod', etc.")
	}

	v := viper.New()

	// Load config file based on APP_ENV
	configName := fmt.Sprintf("config.%s", appEnv)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")

	// Search paths: root, internal/config, and test working directory
	v.AddConfigPath(".")
	v.AddConfigPath("internal/config")
	configDir := filepath.Join("internal", "config")
	v.AddConfigPath(configDir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s/%s.yaml", configDir, configName)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables to struct fields explicitly
	_ = v.BindEnv("db.dsn", "POSTGRE_DB_DSN")
	_ = v.BindEnv("db.host", "DB_HOST")
	_ = v.BindEnv("db.port", "DB_PORT")
	_ = v.BindEnv("db.user", "DB_USER")
	_ = v.BindEnv("db.password", "DB_PASSWORD")
	_ = v.BindEnv("db.name", "DB_NAME")
	_ = v.BindEnv("db.sslmode", "DB_SSLMODE")
	_ = v.BindEnv("jwt.access_secret_key", "JWT_ACCESS_SECRET")
	_ = v.BindEnv("jwt.refresh_secret_key", "JWT_REFRESH_SECRET")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Build GORM DSN if not set but individual values are present
	if cfg.DB.Dsn == "" && cfg.DB.Host != "" {
		cfg.DB.Dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
		)
	}

	return &cfg, nil
}
