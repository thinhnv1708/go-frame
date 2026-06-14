package database

import (
	"fmt"
	"identify/internal/config"
	"identify/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgresDb establishes connection to PostgreSQL and runs GORM auto-migrations.
func ConnectPostgresDb(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(&models.UserModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}

// DisconnectPostgresDb closes the underlying SQL database connection.
func DisconnectPostgresDb(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
