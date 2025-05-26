package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	result := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		result[err.Field()] = fmt.Sprintf("Failed on '%s' tag", err.Tag())
	}
	return result
}
