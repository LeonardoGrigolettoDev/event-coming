package service

import (
	"context"
	"event-coming/internal/config"
	"event-coming/internal/dto"
	"event-coming/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error)
}

// Struct (implementação)
type authServiceImpl struct {
	userRepo  repository.UserRepository
	tokenRepo repository.RefreshTokenRepository // ← Use esta que já existe
	config    *config.JWTConfig
}

// Construtor
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository, // ← Aqui também
	config *config.JWTConfig,
) AuthService { // ← Retorna a interface
	return &authServiceImpl{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    config,
	}
}

// Métodos da implementação
func (s *authServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	return nil, nil
}

func (s *authServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	return nil, nil
}

func (s *authServiceImpl) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	return nil, nil
}
