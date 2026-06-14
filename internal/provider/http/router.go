package http

import (
	"identify/internal/handler"
	"identify/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all routes, middleware, and handlers on the Gin engine.
func SetupRoutes(engine *gin.Engine, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Recover middleware to prevent server crash on panic
	engine.Use(gin.Recovery())

	// Global error handling middleware (from project 1)
	engine.Use(middleware.ErrorMiddleware())

	// Health check endpoint (from project 2)
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register sub-routers
	NewUserRouter(userHandler).Register(engine)
	NewAuthRouter(authHandler).Register(engine)
}
