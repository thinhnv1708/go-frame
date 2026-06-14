package middleware

import (
	"errors"
	"identify/internal/dto/response"
	"identify/internal/exception"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware intercepts errors appended to the Gin context and formats them as standard JSON.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *exception.AppError
		if errors.As(err, &appErr) {
			errResponse := response.NewErrorResponse(appErr.Code, appErr.Message)
			c.JSON(appErr.HttpStatus, errResponse)
			return
		}

		internalErr := exception.INTERNAL_SERVER_ERROR
		c.JSON(internalErr.HttpStatus, response.NewErrorResponse(internalErr.Code, internalErr.Message))
	}
}
