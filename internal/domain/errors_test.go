package domain_test

import (
	"testing"

	"event-coming/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedMsg string
	}{
		{
			name:        "ErrNotFound",
			err:         domain.ErrNotFound,
			expectedMsg: "resource not found",
		},
		{
			name:        "ErrUnauthorized",
			err:         domain.ErrUnauthorized,
			expectedMsg: "unauthorized",
		},
		{
			name:        "ErrForbidden",
			err:         domain.ErrForbidden,
			expectedMsg: "forbidden",
		},
		{
			name:        "ErrConflict",
			err:         domain.ErrConflict,
			expectedMsg: "resource already exists",
		},
		{
			name:        "ErrInvalidInput",
			err:         domain.ErrInvalidInput,
			expectedMsg: "invalid input",
		},
		{
			name:        "ErrInternalServer",
			err:         domain.ErrInternalServer,
			expectedMsg: "internal server error",
		},
		{
			name:        "ErrInvalidCredentials",
			err:         domain.ErrInvalidCredentials,
			expectedMsg: "invalid credentials",
		},
		{
			name:        "ErrTokenExpired",
			err:         domain.ErrTokenExpired,
			expectedMsg: "token expired",
		},
		{
			name:        "ErrInvalidToken",
			err:         domain.ErrInvalidToken,
			expectedMsg: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Equal(t, tt.expectedMsg, tt.err.Error())
		})
	}
}

func TestDomainErrors_AreDistinct(t *testing.T) {
	errors := []error{
		domain.ErrNotFound,
		domain.ErrUnauthorized,
		domain.ErrForbidden,
		domain.ErrConflict,
		domain.ErrInvalidInput,
		domain.ErrInternalServer,
		domain.ErrInvalidCredentials,
		domain.ErrTokenExpired,
		domain.ErrInvalidToken,
	}

	// Verify all errors are distinct
	for i, err1 := range errors {
		for j, err2 := range errors {
			if i != j {
				assert.NotEqual(t, err1, err2, "Errors at index %d and %d should be different", i, j)
			}
		}
	}
}

func TestDomainErrors_NotNil(t *testing.T) {
	assert.NotNil(t, domain.ErrNotFound)
	assert.NotNil(t, domain.ErrUnauthorized)
	assert.NotNil(t, domain.ErrForbidden)
	assert.NotNil(t, domain.ErrConflict)
	assert.NotNil(t, domain.ErrInvalidInput)
	assert.NotNil(t, domain.ErrInternalServer)
	assert.NotNil(t, domain.ErrInvalidCredentials)
	assert.NotNil(t, domain.ErrTokenExpired)
	assert.NotNil(t, domain.ErrInvalidToken)
}
