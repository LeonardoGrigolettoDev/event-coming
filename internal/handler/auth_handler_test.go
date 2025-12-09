package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/dto"
	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockAuthService implements service.AuthService for testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RefreshResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ForgotPasswordResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResetPasswordResponse), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful login",
			body: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("dto.LoginRequest")).Return(&dto.LoginResponse{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, mock.AnythingOfType("dto.LoginRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid json",
			body:           "invalid json",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful registration",
			body: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("dto.RegisterRequest")).Return(&dto.RegisterResponse{
					ID:    "user-id",
					Name:  "John Doe",
					Email: "john@example.com",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "registration fails",
			body: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Register", mock.Anything, mock.AnythingOfType("dto.RegisterRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful refresh",
			body: dto.RefreshRequest{
				RefreshToken: "valid-refresh-token",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Refresh", mock.Anything, mock.AnythingOfType("dto.RefreshRequest")).Return(&dto.RefreshResponse{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
					ExpiresIn:    3600,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid token",
			body: dto.RefreshRequest{
				RefreshToken: "invalid-token",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Refresh", mock.Anything, mock.AnythingOfType("dto.RefreshRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Refresh(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful logout",
			body: dto.LogoutRequest{
				RefreshToken: "valid-refresh-token",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Logout", mock.Anything, mock.AnythingOfType("dto.LogoutRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "logout fails",
			body: dto.LogoutRequest{
				RefreshToken: "invalid-token",
			},
			setupMock: func(m *MockAuthService) {
				m.On("Logout", mock.Anything, mock.AnythingOfType("dto.LogoutRequest")).Return(assert.AnError)
			},
			expectedStatus: http.StatusBadRequest, // Handler returns 400 for invalid token
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Logout(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful forgot password",
			body: dto.ForgotPasswordRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ForgotPassword", mock.Anything, mock.AnythingOfType("dto.ForgotPasswordRequest")).Return(&dto.ForgotPasswordResponse{
					Message: "Email sent",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "forgot password fails - returns 200 for security",
			body: dto.ForgotPasswordRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ForgotPassword", mock.Anything, mock.AnythingOfType("dto.ForgotPasswordRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusOK, // Handler returns 200 even on error for security reasons
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.ForgotPassword(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "successful reset password",
			body: dto.ResetPasswordRequest{
				Token:       "valid-token",
				NewPassword: "newpassword123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, mock.AnythingOfType("dto.ResetPasswordRequest")).Return(&dto.ResetPasswordResponse{
					Message: "Password reset successful",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "reset password fails with invalid token",
			body: dto.ResetPasswordRequest{
				Token:       "invalid-token",
				NewPassword: "newpassword123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, mock.AnythingOfType("dto.ResetPasswordRequest")).Return(nil, service.ErrInvalidToken)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "reset password fails with generic error",
			body: dto.ResetPasswordRequest{
				Token:       "some-token",
				NewPassword: "newpassword123",
			},
			setupMock: func(m *MockAuthService) {
				m.On("ResetPassword", mock.Anything, mock.AnythingOfType("dto.ResetPasswordRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "invalid json",
			body:           "invalid",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			handler := NewAuthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.ResetPassword(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
