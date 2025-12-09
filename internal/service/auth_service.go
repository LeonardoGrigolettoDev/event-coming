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
	Logout(ctx context.Context, req dto.LogoutRequest) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

type authServiceImpl struct {
	userRepo          repository.UserRepository
	tokenRepo         repository.RefreshTokenRepository
	passwordResetRepo repository.PasswordResetTokenRepository
	entityRepo        repository.EntityRepository
	config            *config.JWTConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	passwordResetRepo repository.PasswordResetTokenRepository,
	entityRepo repository.EntityRepository,
	config *config.JWTConfig,
) AuthService {
	return &authServiceImpl{
		userRepo:          userRepo,
		tokenRepo:         tokenRepo,
		passwordResetRepo: passwordResetRepo,
		entityRepo:        entityRepo,
		config:            config,
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
			Active:           true,
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
		"sub":     user.ID.String(),
		"user_id": user.ID.String(),
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(s.config.AccessExpiresIn).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Get user's primary entity and role (first entity association)
	userEntities, err := s.userRepo.GetUserEntities(context.Background(), user.ID)
	if err == nil && len(userEntities) > 0 {
		// Use the first entity as the primary one
		primaryEntity := userEntities[0]
		claims["entity_id"] = primaryEntity.EntityID.String()
		claims["role"] = string(primaryEntity.Role)
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

// ==================== LOGOUT ====================

func (s *authServiceImpl) Logout(ctx context.Context, req dto.LogoutRequest) error {
	// Hash do token para buscar no banco
	tokenHash := s.hashToken(req.RefreshToken)

	// Revogar o refresh token
	if err := s.tokenRepo.RevokeByToken(ctx, tokenHash); err != nil {
		return ErrInvalidToken
	}

	return nil
}

// ==================== FORGOT PASSWORD ====================

func (s *authServiceImpl) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	// 1. Buscar usuário pelo email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		// Retorna sucesso mesmo se usuário não existir (segurança)
		return &dto.ForgotPasswordResponse{
			Message: "If an account with this email exists, a password reset link has been sent.",
		}, nil
	}

	// 2. Deletar tokens antigos do usuário
	_ = s.passwordResetRepo.DeleteByUserID(ctx, user.ID)

	// 3. Gerar novo token
	rawToken := uuid.New().String()
	tokenHash := s.hashToken(rawToken)

	resetToken := &domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token válido por 1 hora
		CreatedAt: time.Now(),
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return nil, err
	}

	// 4. TODO: Enviar email com o token
	// Por enquanto, apenas logamos o token (em produção, enviar por email/WhatsApp)
	// O token a ser enviado é: rawToken
	// A URL seria algo como: https://app.example.com/reset-password?token=rawToken

	return &dto.ForgotPasswordResponse{
		Message: "If an account with this email exists, a password reset link has been sent.",
	}, nil
}

// ==================== RESET PASSWORD ====================

func (s *authServiceImpl) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	// 1. Hash do token
	tokenHash := s.hashToken(req.Token)

	// 2. Buscar token no banco
	resetToken, err := s.passwordResetRepo.GetByToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 3. Buscar usuário
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// 4. Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 5. Atualizar senha do usuário
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 6. Marcar token como usado
	_ = s.passwordResetRepo.MarkAsUsed(ctx, resetToken.ID)

	// 7. Revogar todos os refresh tokens do usuário (força re-login)
	_ = s.tokenRepo.RevokeAllByUserID(ctx, user.ID)

	return &dto.ResetPasswordResponse{
		Message: "Password has been reset successfully. Please login with your new password.",
	}, nil
}
