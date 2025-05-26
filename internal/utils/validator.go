package utils

import (
	"github.com/XSidik/go-fiber-template/internal/models"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type LocalValidator struct {
	models.XValidator
}

func (v LocalValidator) Validate(data interface{}) []models.ErrorResponse {
	validationErrors := []models.ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem models.ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
