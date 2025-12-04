package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is a singleton validator instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors formats validator errors into a slice of ValidationError
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: formatErrorMessage(e),
			})
		}
	}

	return errors
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param())
	case "e164":
		return "Invalid phone number format (E.164 required)"
	case "latitude":
		return "Invalid latitude value"
	case "longitude":
		return "Invalid longitude value"
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s validation failed on %s", e.Field(), e.Tag())
	}
}
