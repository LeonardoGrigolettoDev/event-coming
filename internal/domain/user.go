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
	ID            uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string     `json:"email" db:"email" gorm:"size:255;uniqueIndex;not null"`
	PasswordHash  string     `json:"-" db:"password_hash" gorm:"size:255;not null"`
	Name          string     `json:"name" db:"name" gorm:"size:100;not null"`
	Phone         *string    `json:"phone_number,omitempty" db:"phone_number" gorm:"size:20"`
	Active        bool       `json:"active" db:"active" gorm:"default:true;not null"`
	EmailVerified bool       `json:"email_verified" db:"email_verified" gorm:"default:false;not null"`
	PhoneVerified bool       `json:"phone_verified" db:"phone_verified" gorm:"default:false;not null"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// UserOrganization represents a user's membership in an organization
type UserOrganization struct {
	ID             uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id" gorm:"type:uuid;not null;index"`
	Role           UserRole  `json:"role" db:"role" gorm:"size:50;not null;default:'org_viewer'"`
	CreatedAt      time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

func (UserOrganization) TableName() string {
	return "user_organizations"
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
	ID        uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index"`
	Token     string     `json:"token" db:"token" gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at" gorm:"index"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
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
	UserID         uuid.UUID  `json:"user_id"`
	Email          string     `json:"email"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Role           *UserRole  `json:"role,omitempty"`
}
