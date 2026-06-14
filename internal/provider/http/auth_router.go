package http

import (
	"identify/internal/handler"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	handler *handler.AuthHandler
}

// NewAuthRouter creates a new AuthRouter instance.
func NewAuthRouter(h *handler.AuthHandler) *AuthRouter {
	return &AuthRouter{handler: h}
}

// Register registers Auth routes on the Gin engine.
func (r *AuthRouter) Register(engine *gin.Engine) {
	authGroup := engine.Group("/auth")
	{
		authGroup.POST("/login", r.handler.Login)
	}
}
