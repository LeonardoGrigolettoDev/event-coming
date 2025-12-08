package domain_test

import (
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_TableName(t *testing.T) {
	user := domain.User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUserEntity_TableName(t *testing.T) {
	userEntity := domain.UserEntity{}
	assert.Equal(t, "user_entities", userEntity.TableName())
}

func TestRefreshToken_TableName(t *testing.T) {
	token := domain.RefreshToken{}
	assert.Equal(t, "refresh_tokens", token.TableName())
}

func TestPasswordResetToken_TableName(t *testing.T) {
	token := domain.PasswordResetToken{}
	assert.Equal(t, "password_reset_tokens", token.TableName())
}

func TestUserRole_Constants(t *testing.T) {
	assert.Equal(t, domain.UserRole("super_admin"), domain.UserRoleSuperAdmin)
	assert.Equal(t, domain.UserRole("entity_owner"), domain.UserRoleEntityOwner)
	assert.Equal(t, domain.UserRole("entity_admin"), domain.UserRoleEntityAdmin)
	assert.Equal(t, domain.UserRole("entity_manager"), domain.UserRoleEntityManager)
	assert.Equal(t, domain.UserRole("entity_viewer"), domain.UserRoleEntityViewer)
}

func TestUser_Fields(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	phone := "+5511999999999"
	lastLogin := now.Add(-1 * time.Hour)

	user := domain.User{
		ID:            id,
		Email:         "test@example.com",
		PasswordHash:  "hashedpassword123",
		Name:          "Test User",
		Phone:         &phone,
		Active:        true,
		EmailVerified: true,
		PhoneVerified: false,
		LastLoginAt:   &lastLogin,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	assert.Equal(t, id, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword123", user.PasswordHash)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, &phone, user.Phone)
	assert.True(t, user.Active)
	assert.True(t, user.EmailVerified)
	assert.False(t, user.PhoneVerified)
	assert.Equal(t, &lastLogin, user.LastLoginAt)
}

func TestUser_NilOptionalFields(t *testing.T) {
	user := domain.User{
		ID:           uuid.New(),
		Email:        "minimal@example.com",
		PasswordHash: "hash",
		Name:         "Minimal User",
		Active:       true,
	}

	assert.Nil(t, user.Phone)
	assert.Nil(t, user.LastLoginAt)
}

func TestUserEntity_Fields(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()
	now := time.Now()

	userEntity := domain.UserEntity{
		ID:        id,
		UserID:    userID,
		EntityID:  entityID,
		Role:      domain.UserRoleEntityAdmin,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, id, userEntity.ID)
	assert.Equal(t, userID, userEntity.UserID)
	assert.Equal(t, entityID, userEntity.EntityID)
	assert.Equal(t, domain.UserRoleEntityAdmin, userEntity.Role)
}

func TestRefreshToken_Fields(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour)
	revokedAt := now.Add(1 * time.Hour)

	token := domain.RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     "token123abc",
		ExpiresAt: expiresAt,
		CreatedAt: now,
		RevokedAt: &revokedAt,
	}

	assert.Equal(t, id, token.ID)
	assert.Equal(t, userID, token.UserID)
	assert.Equal(t, "token123abc", token.Token)
	assert.Equal(t, expiresAt, token.ExpiresAt)
	assert.Equal(t, now, token.CreatedAt)
	assert.Equal(t, &revokedAt, token.RevokedAt)
}

func TestRefreshToken_NotRevoked(t *testing.T) {
	token := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "validtoken",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	assert.Nil(t, token.RevokedAt)
}

func TestPasswordResetToken_Fields(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour)
	usedAt := now.Add(30 * time.Minute)

	token := domain.PasswordResetToken{
		ID:        id,
		UserID:    userID,
		Token:     "resettoken123",
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UsedAt:    &usedAt,
	}

	assert.Equal(t, id, token.ID)
	assert.Equal(t, userID, token.UserID)
	assert.Equal(t, "resettoken123", token.Token)
	assert.Equal(t, expiresAt, token.ExpiresAt)
	assert.Equal(t, now, token.CreatedAt)
	assert.Equal(t, &usedAt, token.UsedAt)
}

func TestPasswordResetToken_NotUsed(t *testing.T) {
	token := domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "unusedtoken",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	assert.Nil(t, token.UsedAt)
}

func TestCreateUserInput(t *testing.T) {
	input := domain.CreateUserInput{
		Email:       "newuser@example.com",
		Password:    "securepassword123",
		Name:        "New User",
		PhoneNumber: "+5511999999999",
	}

	assert.Equal(t, "newuser@example.com", input.Email)
	assert.Equal(t, "securepassword123", input.Password)
	assert.Equal(t, "New User", input.Name)
	assert.Equal(t, "+5511999999999", input.PhoneNumber)
}

func TestLoginInput(t *testing.T) {
	input := domain.LoginInput{
		Email:    "user@example.com",
		Password: "password123",
	}

	assert.Equal(t, "user@example.com", input.Email)
	assert.Equal(t, "password123", input.Password)
}

func TestAuthTokens(t *testing.T) {
	tokens := domain.AuthTokens{
		AccessToken:  "access.token.here",
		RefreshToken: "refresh.token.here",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	assert.Equal(t, "access.token.here", tokens.AccessToken)
	assert.Equal(t, "refresh.token.here", tokens.RefreshToken)
	assert.Equal(t, "Bearer", tokens.TokenType)
	assert.Equal(t, int64(3600), tokens.ExpiresIn)
}

func TestJWTClaims(t *testing.T) {
	userID := uuid.New()
	entityID := uuid.New()
	role := domain.UserRoleEntityAdmin

	claims := domain.JWTClaims{
		UserID:   userID,
		Email:    "user@example.com",
		EntityID: &entityID,
		Role:     &role,
	}

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, &entityID, claims.EntityID)
	assert.Equal(t, &role, claims.Role)
}

func TestJWTClaims_NilOptionalFields(t *testing.T) {
	claims := domain.JWTClaims{
		UserID: uuid.New(),
		Email:  "user@example.com",
	}

	assert.Nil(t, claims.EntityID)
	assert.Nil(t, claims.Role)
}
