package service

import (
	"context"
	"identify/internal/dto/request"
	"identify/internal/dto/response"
)

type UserService interface {
	CreateUser(ctx context.Context, request request.CreateUserRequest) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, updateUserReq request.UpdateUserRequest) (*response.UserResponse, error)
	GetUsers(ctx context.Context) ([]response.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*response.UserResponse, error)
}

type AuthService interface {
	Login(ctx context.Context, request request.LoginRequest) (*response.LoginResponse, error)
}
