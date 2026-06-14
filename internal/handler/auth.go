package handler

import (
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/service"
	"identify/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var reqBody request.LoginRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		validationErr := validation.ParseValidationError(err)
		c.Error(validationErr)
		return
	}

	loginResponse, err := h.authService.Login(c.Request.Context(), reqBody)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(loginResponse))
}
