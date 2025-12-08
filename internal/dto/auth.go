package dto

import "time"

// ==================== LOGIN ====================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // segundos até expirar
}

// ==================== REGISTER ====================

// EntityInput representa os dados opcionais para criar uma entidade junto com o usuário
type EntityInput struct {
	Type        string                 `json:"type" binding:"required,oneof=individual company"`
	Name        string                 `json:"name" binding:"required,min=2,max=200"`
	Email       *string                `json:"email,omitempty" binding:"omitempty,email"`
	PhoneNumber *string                `json:"phone_number,omitempty" binding:"omitempty,max=20"`
	Document    *string                `json:"document,omitempty" binding:"omitempty,max=50"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type RegisterRequest struct {
	Name     string       `json:"name" binding:"required,min=2,max=100"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=8"`
	Phone    string       `json:"phone,omitempty" binding:"omitempty,e164"` // formato: +5511999999999
	Entity   *EntityInput `json:"entity,omitempty"`                         // Entidade opcional a ser criada
}

type RegisterResponse struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Email  string          `json:"email"`
	Entity *EntityResponse `json:"entity,omitempty"` // Entidade criada, se houver
}

// ==================== REFRESH ====================

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // se rotacionar
	ExpiresIn    int64  `json:"expires_in"`
}

// ==================== FORGOT PASSWORD ====================

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"` // mensagem genérica
}

// ==================== RESET PASSWORD ====================

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// ==================== USER INFO (opcional) ====================

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== LOGOUT ====================

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
