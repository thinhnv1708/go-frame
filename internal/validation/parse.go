package validation

import (
	"errors"
	"fmt"
	"identify/internal/exception"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ValidationErrorMessages = map[string]string{
	"min":      "{field} must be at least {param}",
	"max":      "{field} must be at most {param}",
	"required": "{field} is required",
}

func ParseValidationError(err error) *exception.AppError {
	var validationErrs validator.ValidationErrors
	errorMessages := make([]string, 0)

	if errors.As(err, &validationErrs) {
		for _, fieldErr := range validationErrs {
			errorMessage, ok := ValidationErrorMessages[fieldErr.Tag()]
			if !ok {
				errorMessage = fmt.Sprintf("{field} validation failed on '%s'", fieldErr.Tag())
			}

			errorMessage = strings.Replace(errorMessage, "{field}", fieldErr.Field(), 1)
			errorMessage = strings.Replace(errorMessage, "{param}", fieldErr.Param(), 1)

			errorMessages = append(errorMessages, errorMessage)
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	finalErrorMessage := strings.Join(errorMessages, " / ")
	appErr := *exception.VALIDATION_ERROR
	appErr.Message = finalErrorMessage
	return &appErr
}
