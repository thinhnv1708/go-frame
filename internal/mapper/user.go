package mapper

import (
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/entity"
	"identify/internal/exception"
	"time"
)

func ToUserFromCreateReq(request request.CreateUserRequest) (*entity.User, error) {
	dob, err := time.Parse("2006-01-02", request.Dob)
	if err != nil {
		return nil, exception.DOB_PARSE_FAILED
	}

	return &entity.User{
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
		Dob:      dob,
	}, nil
}

func ToUserFromUpdateReq(oldUser *entity.User, request request.UpdateUserRequest) (*entity.User, error) {
	dob, err := time.Parse("2006-01-02", request.Dob)
	if err != nil {
		return nil, exception.DOB_PARSE_FAILED
	}

	return &entity.User{
		ID:       oldUser.ID,
		Name:     request.Name,
		Username: oldUser.Username,
		Password: request.Password,
		Dob:      dob,
	}, nil
}

func ToUserResponse(user *entity.User) *response.UserResponse {
	dob := user.Dob.Format("2006-01-02")

	return &response.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Dob:      dob,
	}
}
