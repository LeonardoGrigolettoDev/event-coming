package postgres

import (
	"context"
	"errors"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken

	result := r.db.WithContext(ctx).
		Where("token = ? AND revoked_at IS NULL AND expires_at > ?", token, time.Now()).
		First(&refreshToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &refreshToken, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)

	if result.Error != nil {
		return result.Error
	}

	// Não retorna erro se não encontrou nenhum token
	return nil
}

func (r *refreshTokenRepository) RevokeByToken(ctx context.Context, tokenHash string) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("token = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", now)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&domain.RefreshToken{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
