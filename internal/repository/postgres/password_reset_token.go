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

type passwordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository creates a new password reset token repository
func NewPasswordResetTokenRepository(db *gorm.DB) repository.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(token)
	return result.Error
}

func (r *passwordResetTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	var token domain.PasswordResetToken

	result := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ? AND used_at IS NULL", tokenHash, time.Now()).
		First(&token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &token, nil
}

func (r *passwordResetTokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&domain.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", now)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *passwordResetTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.PasswordResetToken{})

	return result.Error
}

func (r *passwordResetTokenRepository) DeleteExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Where("expires_at < ? OR used_at IS NOT NULL", time.Now()).
		Delete(&domain.PasswordResetToken{})

	return result.Error
}
