package repository

import (
	"context"
	"identify/internal/entity"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUsers(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	FindUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type TransactionManager interface {
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}
