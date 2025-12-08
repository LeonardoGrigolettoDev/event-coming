package service

import (
	"context"
	"testing"
	"time"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/testutil/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func newTestJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		AccessSecret:     "test-secret-key-for-jwt-testing",
		RefreshSecret:    "test-refresh-secret-key",
		AccessExpiresIn:  15 * time.Minute,
		RefreshExpiresIn: 7 * 24 * time.Hour,
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.RegisterRequest
		setupMock func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository, *mocks.MockPasswordResetTokenRepository, *mocks.MockEntityRepository)
		wantErr   error
	}{
		{
			name: "successful registration without entity",
			req: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "+5511999999999",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository, entityRepo *mocks.MockEntityRepository) {
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			req: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository, entityRepo *mocks.MockEntityRepository) {
				existingUser := &domain.User{ID: uuid.New(), Email: "john@example.com"}
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(existingUser, nil)
			},
			wantErr: ErrEmailAlreadyExists,
		},
		{
			name: "registration with entity",
			req: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Entity: &dto.EntityInput{
					Type: "individual",
					Name: "John's Company",
				},
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository, entityRepo *mocks.MockEntityRepository) {
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				entityRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Entity")).Return(nil)
				userRepo.On("AddToEntity", mock.Anything, mock.AnythingOfType("*domain.UserEntity")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "registration with entity - document conflict",
			req: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Entity: &dto.EntityInput{
					Type:     "individual",
					Name:     "John's Company",
					Document: strPtr("12345678901"),
				},
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository, entityRepo *mocks.MockEntityRepository) {
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				existingEntity := &domain.Entity{ID: uuid.New()}
				entityRepo.On("GetByDocument", mock.Anything, "12345678901").Return(existingEntity, nil)
			},
			wantErr: domain.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(userRepo, tokenRepo, passRepo, entityRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			result, err := svc.Register(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
				assert.Equal(t, tt.req.Email, result.Email)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	validPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		req       dto.LoginRequest
		setupMock func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository)
		wantErr   error
	}{
		{
			name: "successful login",
			req: dto.LoginRequest{
				Email:    "john@example.com",
				Password: validPassword,
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "john@example.com",
					Name:         "John Doe",
					PasswordHash: string(hashedPassword),
					Active:       true,
				}
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
				tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
				userRepo.On("UpdateLastLogin", mock.Anything, user.ID, mock.AnythingOfType("time.Time")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: dto.LoginRequest{
				Email:    "unknown@example.com",
				Password: "password123",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				userRepo.On("GetByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			req: dto.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "john@example.com",
					PasswordHash: string(hashedPassword),
					Active:       true,
				}
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "inactive user",
			req: dto.LoginRequest{
				Email:    "john@example.com",
				Password: validPassword,
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "john@example.com",
					PasswordHash: string(hashedPassword),
					Active:       false,
				}
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(userRepo, tokenRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			result, err := svc.Login(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Greater(t, result.ExpiresIn, int64(0))
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		req       dto.RefreshRequest
		setupMock func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository)
		wantErr   error
	}{
		{
			name: "successful refresh",
			req:  dto.RefreshRequest{RefreshToken: "valid-token"},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				storedToken := &domain.RefreshToken{
					ID:        uuid.New(),
					UserID:    userID,
					ExpiresAt: time.Now().Add(time.Hour),
				}
				tokenRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(storedToken, nil)
				tokenRepo.On("Revoke", mock.Anything, storedToken.ID).Return(nil)

				user := &domain.User{
					ID:    userID,
					Email: "john@example.com",
					Name:  "John Doe",
				}
				userRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
				tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "token not found",
			req:  dto.RefreshRequest{RefreshToken: "invalid-token"},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				tokenRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "token revoked",
			req:  dto.RefreshRequest{RefreshToken: "revoked-token"},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				now := time.Now()
				storedToken := &domain.RefreshToken{
					ID:        uuid.New(),
					UserID:    userID,
					ExpiresAt: time.Now().Add(time.Hour),
					RevokedAt: &now,
				}
				tokenRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(storedToken, nil)
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "token expired",
			req:  dto.RefreshRequest{RefreshToken: "expired-token"},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository) {
				storedToken := &domain.RefreshToken{
					ID:        uuid.New(),
					UserID:    userID,
					ExpiresAt: time.Now().Add(-time.Hour), // expired
				}
				tokenRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(storedToken, nil)
			},
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(userRepo, tokenRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			result, err := svc.Refresh(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			tokenRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.LogoutRequest
		setupMock func(*mocks.MockRefreshTokenRepository)
		wantErr   error
	}{
		{
			name: "successful logout",
			req:  dto.LogoutRequest{RefreshToken: "valid-token"},
			setupMock: func(tokenRepo *mocks.MockRefreshTokenRepository) {
				tokenRepo.On("RevokeByToken", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "token not found",
			req:  dto.LogoutRequest{RefreshToken: "invalid-token"},
			setupMock: func(tokenRepo *mocks.MockRefreshTokenRepository) {
				tokenRepo.On("RevokeByToken", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrNotFound)
			},
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(tokenRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			err := svc.Logout(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			tokenRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ForgotPassword(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.ForgotPasswordRequest
		setupMock func(*mocks.MockUserRepository, *mocks.MockPasswordResetTokenRepository)
		wantErr   bool
	}{
		{
			name: "user exists",
			req:  dto.ForgotPasswordRequest{Email: "john@example.com"},
			setupMock: func(userRepo *mocks.MockUserRepository, passRepo *mocks.MockPasswordResetTokenRepository) {
				user := &domain.User{ID: uuid.New(), Email: "john@example.com"}
				userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)
				passRepo.On("DeleteByUserID", mock.Anything, user.ID).Return(nil)
				passRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PasswordResetToken")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user not found - still returns success",
			req:  dto.ForgotPasswordRequest{Email: "unknown@example.com"},
			setupMock: func(userRepo *mocks.MockUserRepository, passRepo *mocks.MockPasswordResetTokenRepository) {
				userRepo.On("GetByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(userRepo, passRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			result, err := svc.ForgotPassword(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Message)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ResetPassword(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		req       dto.ResetPasswordRequest
		setupMock func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository, *mocks.MockPasswordResetTokenRepository)
		wantErr   error
	}{
		{
			name: "successful reset",
			req: dto.ResetPasswordRequest{
				Token:       "valid-token",
				NewPassword: "newpassword123",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository) {
				resetToken := &domain.PasswordResetToken{
					ID:        uuid.New(),
					UserID:    userID,
					ExpiresAt: time.Now().Add(time.Hour),
				}
				passRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(resetToken, nil)

				user := &domain.User{
					ID:    userID,
					Email: "john@example.com",
				}
				userRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				passRepo.On("MarkAsUsed", mock.Anything, resetToken.ID).Return(nil)
				tokenRepo.On("RevokeAllByUserID", mock.Anything, userID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid token",
			req: dto.ResetPasswordRequest{
				Token:       "invalid-token",
				NewPassword: "newpassword123",
			},
			setupMock: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockRefreshTokenRepository, passRepo *mocks.MockPasswordResetTokenRepository) {
				passRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound)
			},
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.MockUserRepository)
			tokenRepo := new(mocks.MockRefreshTokenRepository)
			passRepo := new(mocks.MockPasswordResetTokenRepository)
			entityRepo := new(mocks.MockEntityRepository)

			tt.setupMock(userRepo, tokenRepo, passRepo)

			svc := NewAuthService(userRepo, tokenRepo, passRepo, entityRepo, newTestJWTConfig())
			result, err := svc.ResetPassword(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Message)
			}

			passRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}
