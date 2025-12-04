package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents a user's role in an organization
type UserRole string

const (
	UserRoleSuperAdmin  UserRole = "super_admin"
	UserRoleOrgOwner    UserRole = "org_owner"
	UserRoleOrgAdmin    UserRole = "org_admin"
	UserRoleOrgManager  UserRole = "org_manager"
	UserRoleOrgOperator UserRole = "org_operator"
	UserRoleOrgViewer   UserRole = "org_viewer"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	Name           string    `json:"name" db:"name"`
	PhoneNumber    *string   `json:"phone_number,omitempty" db:"phone_number"`
	Active         bool      `json:"active" db:"active"`
	EmailVerified  bool      `json:"email_verified" db:"email_verified"`
	PhoneVerified  bool      `json:"phone_verified" db:"phone_verified"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// UserOrganization represents a user's membership in an organization
type UserOrganization struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Role           UserRole  `json:"role" db:"role"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserInput holds data for user registration
type CreateUserInput struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,e164"`
}

// LoginInput holds data for user login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// AuthTokens holds access and refresh tokens
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTClaims holds JWT token claims
type JWTClaims struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Role           *UserRole  `json:"role,omitempty"`
}
