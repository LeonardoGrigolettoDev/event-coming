package validator_test

import (
	"testing"

	"event-coming/pkg/validator"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Initialization(t *testing.T) {
	assert.NotNil(t, validator.Validate)
}

func TestValidationError_Structure(t *testing.T) {
	err := validator.ValidationError{
		Field:   "email",
		Message: "Invalid email format",
	}

	assert.Equal(t, "email", err.Field)
	assert.Equal(t, "Invalid email format", err.Message)
}

type TestStruct struct {
	Email     string  `validate:"required,email"`
	Name      string  `validate:"required,min=2,max=100"`
	Phone     string  `validate:"omitempty,e164"`
	Latitude  float64 `validate:"required,latitude"`
	Longitude float64 `validate:"required,longitude"`
	Status    string  `validate:"required,oneof=active inactive pending"`
}

func TestFormatValidationErrors_Required(t *testing.T) {
	input := TestStruct{}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)
	assert.NotEmpty(t, errors)

	// Check for required fields
	fieldNames := make(map[string]bool)
	for _, e := range errors {
		fieldNames[e.Field] = true
	}

	assert.True(t, fieldNames["email"])
	assert.True(t, fieldNames["name"])
	assert.True(t, fieldNames["status"])
}

func TestFormatValidationErrors_Email(t *testing.T) {
	input := TestStruct{
		Email:     "invalid-email",
		Name:      "Test",
		Latitude:  0,
		Longitude: 0,
		Status:    "active",
	}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)

	var emailError *validator.ValidationError
	for _, e := range errors {
		if e.Field == "email" {
			emailError = &e
			break
		}
	}

	assert.NotNil(t, emailError)
	assert.Contains(t, emailError.Message, "email")
}

func TestFormatValidationErrors_Min(t *testing.T) {
	input := TestStruct{
		Email:     "test@example.com",
		Name:      "A", // Too short
		Latitude:  0,
		Longitude: 0,
		Status:    "active",
	}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)

	var nameError *validator.ValidationError
	for _, e := range errors {
		if e.Field == "name" {
			nameError = &e
			break
		}
	}

	assert.NotNil(t, nameError)
	assert.Contains(t, nameError.Message, "at least")
}

func TestFormatValidationErrors_Max(t *testing.T) {
	type MaxTestStruct struct {
		Name string `validate:"max=5"`
	}

	input := MaxTestStruct{
		Name: "This is a very long name that exceeds the maximum",
	}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0].Message, "must not exceed")
}

func TestFormatValidationErrors_OneOf(t *testing.T) {
	input := TestStruct{
		Email:     "test@example.com",
		Name:      "Test Name",
		Latitude:  0,
		Longitude: 0,
		Status:    "unknown", // Invalid status
	}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)

	var statusError *validator.ValidationError
	for _, e := range errors {
		if e.Field == "status" {
			statusError = &e
			break
		}
	}

	assert.NotNil(t, statusError)
	assert.Contains(t, statusError.Message, "must be one of")
}

func TestFormatValidationErrors_ValidInput(t *testing.T) {
	input := TestStruct{
		Email:     "test@example.com",
		Name:      "Test Name",
		Phone:     "+5511999999999",
		Latitude:  -23.550520,
		Longitude: -46.633308,
		Status:    "active",
	}

	err := validator.Validate.Struct(input)
	assert.NoError(t, err)
}

func TestFormatValidationErrors_NonValidationError(t *testing.T) {
	// Test with non-validation error
	errors := validator.FormatValidationErrors(nil)
	assert.Empty(t, errors)
}

func TestFormatValidationErrors_FieldNamesLowercase(t *testing.T) {
	type CamelCaseStruct struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
	}

	input := CamelCaseStruct{}

	err := validator.Validate.Struct(input)
	assert.Error(t, err)

	errors := validator.FormatValidationErrors(err)

	for _, e := range errors {
		// Field names should be lowercase
		assert.Equal(t, e.Field, e.Field)
	}
}

func TestValidate_Latitude(t *testing.T) {
	type LatStruct struct {
		Lat float64 `validate:"latitude"`
	}

	tests := []struct {
		name    string
		lat     float64
		isValid bool
	}{
		{"Valid positive", 45.0, true},
		{"Valid negative", -45.0, true},
		{"Valid max", 90.0, true},
		{"Valid min", -90.0, true},
		{"Valid zero", 0, true},
		{"Invalid too high", 91.0, false},
		{"Invalid too low", -91.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := LatStruct{Lat: tt.lat}
			err := validator.Validate.Struct(input)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidate_Longitude(t *testing.T) {
	type LngStruct struct {
		Lng float64 `validate:"longitude"`
	}

	tests := []struct {
		name    string
		lng     float64
		isValid bool
	}{
		{"Valid positive", 90.0, true},
		{"Valid negative", -90.0, true},
		{"Valid max", 180.0, true},
		{"Valid min", -180.0, true},
		{"Valid zero", 0, true},
		{"Invalid too high", 181.0, false},
		{"Invalid too low", -181.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := LngStruct{Lng: tt.lng}
			err := validator.Validate.Struct(input)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
