package di

import (
	"gorm.io/gorm"

	"identify/internal/config"
	"identify/internal/logger"
	"identify/internal/provider"
)

// Container holds all application dependencies.
type Container struct {
	Config   *config.Config
	Logger   logger.Logger
	DB       *gorm.DB
	Registry *provider.Registry
}
