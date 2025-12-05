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

// ==================== USER-ORGANIZATION ====================

func (r *userRepository) AddToOrganization(ctx context.Context, userOrg *domain.UserOrganization) error {
	if userOrg.ID == uuid.Nil {
		userOrg.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(userOrg)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) RemoveFromOrganization(ctx context.Context, userID, orgID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&domain.UserOrganization{}, "user_id = ? AND organization_id = ?", userID, orgID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.UserOrganization, error) {
	var userOrgs []*domain.UserOrganization

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&userOrgs)

	if result.Error != nil {
		return nil, result.Error
	}

	return userOrgs, nil
}

func (r *userRepository) GetOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User

	result := r.db.WithContext(ctx).
		Joins("JOIN user_organizations ON user_organizations.user_id = users.id").
		Where("user_organizations.organization_id = ?", orgID).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
