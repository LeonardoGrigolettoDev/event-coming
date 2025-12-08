package dto_test

import (
	"testing"
	"time"

	"event-coming/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestLoginRequest(t *testing.T) {
	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "securepassword123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "securepassword123", req.Password)
}

func TestLoginResponse(t *testing.T) {
	resp := dto.LoginResponse{
		AccessToken:  "access.token.here",
		RefreshToken: "refresh.token.here",
		ExpiresIn:    3600,
	}

	assert.Equal(t, "access.token.here", resp.AccessToken)
	assert.Equal(t, "refresh.token.here", resp.RefreshToken)
	assert.Equal(t, int64(3600), resp.ExpiresIn)
}

func TestEntityInput(t *testing.T) {
	email := "entity@example.com"
	phone := "+5511999999999"
	doc := "12345678901"

	input := dto.EntityInput{
		Type:        "individual",
		Name:        "Test Entity",
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		Metadata:    map[string]interface{}{"key": "value"},
	}

	assert.Equal(t, "individual", input.Type)
	assert.Equal(t, "Test Entity", input.Name)
	assert.Equal(t, &email, input.Email)
	assert.Equal(t, &phone, input.PhoneNumber)
	assert.Equal(t, &doc, input.Document)
	assert.Equal(t, "value", input.Metadata["key"])
}

func TestRegisterRequest(t *testing.T) {
	email := "entity@example.com"
	entityInput := dto.EntityInput{
		Type:  "company",
		Name:  "Test Company",
		Email: &email,
	}

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "user@example.com",
		Password: "securepassword123",
		Phone:    "+5511999999999",
		Entity:   &entityInput,
	}

	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "user@example.com", req.Email)
	assert.Equal(t, "securepassword123", req.Password)
	assert.Equal(t, "+5511999999999", req.Phone)
	assert.NotNil(t, req.Entity)
	assert.Equal(t, "company", req.Entity.Type)
}

func TestRegisterRequest_WithoutEntity(t *testing.T) {
	req := dto.RegisterRequest{
		Name:     "Simple User",
		Email:    "simple@example.com",
		Password: "password123",
	}

	assert.Nil(t, req.Entity)
}

func TestRegisterResponse(t *testing.T) {
	entityResp := &dto.EntityResponse{
		Name: "Test Entity",
	}

	resp := dto.RegisterResponse{
		ID:     "user-id-123",
		Name:   "Test User",
		Email:  "user@example.com",
		Entity: entityResp,
	}

	assert.Equal(t, "user-id-123", resp.ID)
	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "user@example.com", resp.Email)
	assert.NotNil(t, resp.Entity)
}

func TestRefreshRequest(t *testing.T) {
	req := dto.RefreshRequest{
		RefreshToken: "refresh.token.here",
	}

	assert.Equal(t, "refresh.token.here", req.RefreshToken)
}

func TestRefreshResponse(t *testing.T) {
	resp := dto.RefreshResponse{
		AccessToken:  "new.access.token",
		RefreshToken: "new.refresh.token",
		ExpiresIn:    7200,
	}

	assert.Equal(t, "new.access.token", resp.AccessToken)
	assert.Equal(t, "new.refresh.token", resp.RefreshToken)
	assert.Equal(t, int64(7200), resp.ExpiresIn)
}

func TestForgotPasswordRequest(t *testing.T) {
	req := dto.ForgotPasswordRequest{
		Email: "user@example.com",
	}

	assert.Equal(t, "user@example.com", req.Email)
}

func TestForgotPasswordResponse(t *testing.T) {
	resp := dto.ForgotPasswordResponse{
		Message: "Password reset email sent",
	}

	assert.Equal(t, "Password reset email sent", resp.Message)
}

func TestResetPasswordRequest(t *testing.T) {
	req := dto.ResetPasswordRequest{
		Token:       "reset-token-123",
		NewPassword: "newpassword456",
	}

	assert.Equal(t, "reset-token-123", req.Token)
	assert.Equal(t, "newpassword456", req.NewPassword)
}

func TestResetPasswordResponse(t *testing.T) {
	resp := dto.ResetPasswordResponse{
		Message: "Password reset successfully",
	}

	assert.Equal(t, "Password reset successfully", resp.Message)
}

func TestUserResponse(t *testing.T) {
	now := time.Now()

	resp := dto.UserResponse{
		ID:        "user-id-123",
		Name:      "Test User",
		Email:     "user@example.com",
		Phone:     "+5511999999999",
		CreatedAt: now,
	}

	assert.Equal(t, "user-id-123", resp.ID)
	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "user@example.com", resp.Email)
	assert.Equal(t, "+5511999999999", resp.Phone)
	assert.Equal(t, now, resp.CreatedAt)
}

func TestLogoutRequest(t *testing.T) {
	req := dto.LogoutRequest{
		RefreshToken: "token.to.revoke",
	}

	assert.Equal(t, "token.to.revoke", req.RefreshToken)
}

func TestLogoutResponse(t *testing.T) {
	resp := dto.LogoutResponse{
		Message: "Successfully logged out",
	}

	assert.Equal(t, "Successfully logged out", resp.Message)
}
