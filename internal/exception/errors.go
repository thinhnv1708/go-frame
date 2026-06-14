package exception

import "net/http"

var (
	INTERNAL_SERVER_ERROR = &AppError{
		Code:       1000,
		Message:    "internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
	DOB_PARSE_FAILED = &AppError{
		Code:       1001,
		Message:    "invalid date of birth format, expected YYYY-MM-DD",
		HttpStatus: http.StatusBadRequest,
	}
	VALIDATION_ERROR = &AppError{
		Code:       1004,
		Message:    "",
		HttpStatus: http.StatusBadRequest,
	}
	USER_NOT_FOUND = &AppError{
		Code:       1002,
		Message:    "user not found",
		HttpStatus: http.StatusNotFound,
	}
	PASSWORD_HASHING_FAILED = &AppError{
		Code:       1003,
		Message:    "failed to hash password",
		HttpStatus: http.StatusInternalServerError,
	}
	INVALID_ACCESS_TOKEN = &AppError{
		Code:       1005,
		Message:    "invalid access token",
		HttpStatus: http.StatusUnauthorized,
	}
	INVALID_REFRESH_TOKEN = &AppError{
		Code:       1006,
		Message:    "invalid refresh token",
		HttpStatus: http.StatusUnauthorized,
	}
	LOGIN_FAILED_USERNAME_OR_PASSWORD_INCORRECT = &AppError{
		Code:       1007,
		Message:    "login failed: username or password is incorrect",
		HttpStatus: http.StatusUnauthorized,
	}
	ACCESS_TOKEN_GENERATION_FAILED = &AppError{
		Code:       1008,
		Message:    "internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
	REFRESH_TOKEN_GENERATION_FAILED = &AppError{
		Code:       1009,
		Message:    "internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
)
