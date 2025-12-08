package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Erros do service
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error)
}

type authServiceImpl struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.RefreshTokenRepository
	entityRepo repository.EntityRepository
	config     *config.JWTConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	entityRepo repository.EntityRepository,
	config *config.JWTConfig,
) AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		entityRepo: entityRepo,
		config:     config,
	}
}

// ==================== REGISTER ====================

func (s *authServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 1. Verificar se email já existe
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// 2. Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Criar usuário
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        &req.Phone,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 4. Criar entidade se fornecida
	var entityResponse *dto.EntityResponse
	if req.Entity != nil {
		// Verificar documento único se fornecido
		if req.Entity.Document != nil && *req.Entity.Document != "" {
			existing, _ := s.entityRepo.GetByDocument(ctx, *req.Entity.Document)
			if existing != nil {
				return nil, domain.ErrConflict
			}
		}

		entity := &domain.Entity{
			ID:               uuid.New(),
			Type:             domain.EntityType(req.Entity.Type),
			Name:             req.Entity.Name,
			Email:            req.Entity.Email,
			PhoneNumber:      req.Entity.PhoneNumber,
			Document:         req.Entity.Document,
			IsActive:         true,
			EntityPermission: domain.EntityPermissionAdmin, // Criador é admin
			Metadata:         req.Entity.Metadata,
		}

		if err := s.entityRepo.Create(ctx, entity); err != nil {
			return nil, err
		}

		// 5. Associar usuário à entidade como owner
		userEntity := &domain.UserEntity{
			ID:       uuid.New(),
			UserID:   user.ID,
			EntityID: entity.ID,
			Role:     domain.UserRoleEntityOwner,
		}

		if err := s.userRepo.AddToEntity(ctx, userEntity); err != nil {
			return nil, err
		}

		entityResponse = dto.ToEntityResponse(entity)
	}

	// 6. Retornar resposta (sem tokens - usuário precisa fazer login)
	return &dto.RegisterResponse{
		ID:     user.ID.String(),
		Name:   user.Name,
		Email:  user.Email,
		Entity: entityResponse,
	}, nil
}

// ==================== LOGIN ====================

func (s *authServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Buscar usuário por email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Verificar se usuário está ativo
	if !user.Active {
		return nil, ErrInvalidCredentials
	}

	// 4. Gerar tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	// 5. Atualizar último login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now())

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiresIn.Seconds()),
	}, nil
}

// ==================== REFRESH ====================

func (s *authServiceImpl) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	// 1. Hash do token recebido
	tokenHash := s.hashToken(req.RefreshToken)

	// 2. Buscar token no banco
	storedToken, err := s.tokenRepo.GetByToken(ctx, tokenHash)
	if err != nil || storedToken == nil {
		return nil, ErrInvalidToken
	}

	// 3. Verificar se não foi revogado
	if storedToken.RevokedAt != nil {
		return nil, ErrInvalidToken
	}

	// 4. Verificar se não expirou
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// 5. Buscar usuário
	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// 6. Revogar token antigo
	_ = s.tokenRepo.Revoke(ctx, storedToken.ID)

	// 7. Gerar novos tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.AccessExpiresIn.Seconds()),
	}, nil
}

// ==================== HELPERS ====================

func (s *authServiceImpl) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(s.config.AccessExpiresIn).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.AccessSecret))
}

func (s *authServiceImpl) generateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	// 1. Gerar token aleatório
	rawToken := uuid.New().String()
	tokenHash := s.hashToken(rawToken)

	// 2. Salvar no banco
	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokenHash, // Salvamos o hash, não o token
		ExpiresAt: time.Now().Add(s.config.RefreshExpiresIn),
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	// 3. Retornar token raw (não o hash)
	return rawToken, nil
}

func (s *authServiceImpl) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
