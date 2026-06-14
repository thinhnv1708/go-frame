package repository

import (
	"context"
	"errors"
	"identify/internal/database/models"
	"identify/internal/entity"

	"gorm.io/gorm"
)

type UserRepositoryImp struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance using GORM.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImp{
		db: db,
	}
}

func (r *UserRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	return GetDBFromContext(ctx, r.db)
}

func toEntity(user *models.UserModel) *entity.User {
	if user == nil {
		return nil
	}
	return &entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Dob:      user.Dob,
	}
}

func toModel(user *entity.User) *models.UserModel {
	if user == nil {
		return nil
	}
	return &models.UserModel{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Dob:      user.Dob,
	}
}

func (r *UserRepositoryImp) SaveUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	userDoc := toModel(user)

	if err := r.getDB(ctx).Create(userDoc).Error; err != nil {
		return nil, err
	}

	return toEntity(userDoc), nil
}

func (r *UserRepositoryImp) FindUsers(ctx context.Context) ([]entity.User, error) {
	var userModels []models.UserModel
	if err := r.getDB(ctx).Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := make([]entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = *toEntity(&userModel)
	}
	return users, nil
}

func (r *UserRepositoryImp) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	userDoc := toModel(user)

	if err := r.getDB(ctx).Save(userDoc).Error; err != nil {
		return nil, err
	}

	return toEntity(userDoc), nil
}

func (r *UserRepositoryImp) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	var user models.UserModel

	err := r.getDB(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toEntity(&user), nil
}

func (r *UserRepositoryImp) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user models.UserModel

	err := r.getDB(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toEntity(&user), nil
}
