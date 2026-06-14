package service

import (
	"context"
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/exception"
	"identify/internal/repository"
	"identify/internal/security"
)

type AuthServiceImp struct {
	userRepository  repository.UserRepository
	passwordEncoder security.PasswordEncoder
	jwtProvider     security.JwtProvider
}

// NewAuthService creates a new AuthService implementation.
func NewAuthService(userRepository repository.UserRepository, passwordEncoder security.PasswordEncoder, jwtProvider security.JwtProvider) AuthService {
	return &AuthServiceImp{
		userRepository:  userRepository,
		passwordEncoder: passwordEncoder,
		jwtProvider:     jwtProvider,
	}
}

func (s *AuthServiceImp) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepository.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, exception.LOGIN_FAILED_USERNAME_OR_PASSWORD_INCORRECT
	}

	verified := s.passwordEncoder.Verify(user.Password, req.Password)
	if !verified {
		return nil, exception.LOGIN_FAILED_USERNAME_OR_PASSWORD_INCORRECT
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, exception.ACCESS_TOKEN_GENERATION_FAILED
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, exception.REFRESH_TOKEN_GENERATION_FAILED
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
