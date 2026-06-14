package service

import (
	"context"
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/exception"
	"identify/internal/mapper"
	"identify/internal/repository"
	"identify/internal/security"

	"github.com/google/uuid"
)

type UserServiceImp struct {
	userRepository  repository.UserRepository
	passwordEncoder security.PasswordEncoder
	txManager       repository.TransactionManager
}

// NewUserService creates a new UserService implementation.
func NewUserService(userRepository repository.UserRepository, passwordEncoder security.PasswordEncoder, txManager repository.TransactionManager) UserService {
	return &UserServiceImp{
		userRepository:  userRepository,
		passwordEncoder: passwordEncoder,
		txManager:       txManager,
	}
}

func (s *UserServiceImp) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	var resp *response.UserResponse
	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		user, err := mapper.ToUserFromCreateReq(req)
		if err != nil {
			return err
		}

		hashedPassword, err := s.passwordEncoder.Encode(user.Password)
		if err != nil {
			return exception.PASSWORD_HASHING_FAILED
		}
		user.Password = hashedPassword

		id := uuid.New()
		user.ID = id.String()

		savedUser, err := s.userRepository.SaveUser(txCtx, user)
		if err != nil {
			return err
		}

		resp = mapper.ToUserResponse(savedUser)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserServiceImp) GetUsers(ctx context.Context) ([]response.UserResponse, error) {
	users, err := s.userRepository.FindUsers(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *mapper.ToUserResponse(&user)
	}

	return userResponses, nil
}

func (s *UserServiceImp) UpdateUser(ctx context.Context, userID string, req request.UpdateUserRequest) (*response.UserResponse, error) {
	var resp *response.UserResponse
	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.userRepository.FindUserByID(txCtx, userID)
		if err != nil {
			return err
		}

		if user == nil {
			return exception.USER_NOT_FOUND
		}

		newUser, err := mapper.ToUserFromUpdateReq(user, req)
		if err != nil {
			return err
		}

		hashedPassword, err := s.passwordEncoder.Encode(newUser.Password)
		if err != nil {
			return exception.PASSWORD_HASHING_FAILED
		}
		newUser.Password = hashedPassword

		updatedUser, err := s.userRepository.UpdateUser(txCtx, newUser)
		if err != nil {
			return err
		}

		resp = mapper.ToUserResponse(updatedUser)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserServiceImp) GetUserByID(ctx context.Context, id string) (*response.UserResponse, error) {
	user, err := s.userRepository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, exception.USER_NOT_FOUND
	}

	return mapper.ToUserResponse(user), nil
}
