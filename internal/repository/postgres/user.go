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

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// ==================== USER CRUD ====================

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("last_login_at", loginTime)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// ==================== USER-ENTITY ====================

func (r *userRepository) AddToEntity(ctx context.Context, userEnt *domain.UserEntity) error {
	if userEnt.ID == uuid.Nil {
		userEnt.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(userEnt)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) RemoveFromEntity(ctx context.Context, userID, entID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&domain.UserEntity{}, "user_id = ? AND entity_id = ?", userID, entID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) GetUserEntities(ctx context.Context, userID uuid.UUID) ([]*domain.UserEntity, error) {
	var userOrgs []*domain.UserEntity

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&userOrgs)

	if result.Error != nil {
		return nil, result.Error
	}

	return userOrgs, nil
}

func (r *userRepository) GetEntityUsers(ctx context.Context, entID uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User

	result := r.db.WithContext(ctx).
		Joins("JOIN user_entities ON user_entities.user_id = users.id").
		Where("user_entities.entity_id = ?", entID).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
