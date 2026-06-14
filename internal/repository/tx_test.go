package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"identify/internal/database/models"
)

func TestGormTransactionManager_Execute(t *testing.T) {
	// Initialize in-memory SQLite for testing GORM transactions
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.UserModel{})
	require.NoError(t, err)

	txManager := NewTransactionManager(db)

	t.Run("Success Commit", func(t *testing.T) {
		ctx := context.Background()
		err := txManager.Execute(ctx, func(txCtx context.Context) error {
			// Retrieve the transaction GORM db from context
			txDB := GetDBFromContext(txCtx, db)
			assert.NotEqual(t, db, txDB) // should be a transaction session

			userModel := &models.UserModel{
				ID:       "test-id-1",
				Name:     "Test User 1",
				Username: "testuser1",
				Password: "password123",
				Dob:      time.Now(),
			}
			return txDB.Create(userModel).Error
		})
		require.NoError(t, err)

		// Verify that it is persisted
		var user models.UserModel
		err = db.First(&user, "id = ?", "test-id-1").Error
		require.NoError(t, err)
		assert.Equal(t, "Test User 1", user.Name)
	})

	t.Run("Failure Rollback", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("something went wrong, rollback please")

		err := txManager.Execute(ctx, func(txCtx context.Context) error {
			txDB := GetDBFromContext(txCtx, db)

			userModel := &models.UserModel{
				ID:       "test-id-2",
				Name:     "Test User 2",
				Username: "testuser2",
				Password: "password123",
				Dob:      time.Now(),
			}
			err := txDB.Create(userModel).Error
			require.NoError(t, err)

			// Return error to trigger rollback
			return expectedErr
		})
		require.ErrorIs(t, err, expectedErr)

		// Verify that it is NOT persisted due to rollback
		var user models.UserModel
		err = db.First(&user, "id = ?", "test-id-2").Error
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("GetDBFromContext returns default if no transaction in context", func(t *testing.T) {
		ctx := context.Background()
		returnedDB := GetDBFromContext(ctx, db)
		assert.NotNil(t, returnedDB)
		var count int64
		err := returnedDB.Model(&models.UserModel{}).Count(&count).Error
		assert.NoError(t, err)
	})
}
