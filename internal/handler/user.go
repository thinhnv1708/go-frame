package handler

import (
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/service"
	"identify/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var body request.CreateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		validationErr := validation.ParseValidationError(err)
		c.Error(validationErr)
		return
	}

	userResponse, err := h.userService.CreateUser(c.Request.Context(), body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(userResponse))
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	userResponses, err := h.userService.GetUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(userResponses))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var body request.UpdateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		validationErr := validation.ParseValidationError(err)
		c.Error(validationErr)
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), userID, body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(userResponse))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(userResponse))
}
