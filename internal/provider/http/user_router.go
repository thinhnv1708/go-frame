package http

import (
	"identify/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	handler *handler.UserHandler
}

// NewUserRouter creates a new UserRouter instance.
func NewUserRouter(h *handler.UserHandler) *UserRouter {
	return &UserRouter{handler: h}
}

// Register registers User routes on the Gin engine.
func (r *UserRouter) Register(engine *gin.Engine) {
	userGroup := engine.Group("/users")
	{
		userGroup.POST("", r.handler.CreateUser)
		userGroup.GET("", r.handler.GetUsers)
		userGroup.GET("/:id", r.handler.GetUser)
		userGroup.PUT("/:id", r.handler.UpdateUser)
	}
}
