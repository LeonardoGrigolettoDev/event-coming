package mocks

import (
	"context"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	args := m.Called(ctx, id, loginTime)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) AddToEntity(ctx context.Context, userEntity *domain.UserEntity) error {
	args := m.Called(ctx, userEntity)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveFromEntity(ctx context.Context, userID, entityID uuid.UUID) error {
	args := m.Called(ctx, userID, entityID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserEntities(ctx context.Context, userID uuid.UUID) ([]*domain.UserEntity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetEntityUsers(ctx context.Context, entityID uuid.UUID) ([]*domain.User, error) {
	args := m.Called(ctx, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockPasswordResetTokenRepository is a mock implementation of PasswordResetTokenRepository
type MockPasswordResetTokenRepository struct {
	mock.Mock
}

func (m *MockPasswordResetTokenRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetTokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockEntityRepository is a mock implementation of EntityRepository
type MockEntityRepository struct {
	mock.Mock
}

func (m *MockEntityRepository) Create(ctx context.Context, entity *domain.Entity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockEntityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entity), args.Error(1)
}

func (m *MockEntityRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateEntityInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockEntityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEntityRepository) List(ctx context.Context, page, perPage int) ([]*domain.Entity, int64, error) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Entity), args.Get(1).(int64), args.Error(2)
}

func (m *MockEntityRepository) ListByParent(ctx context.Context, parentID uuid.UUID, page, perPage int) ([]*domain.Entity, int64, error) {
	args := m.Called(ctx, parentID, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Entity), args.Get(1).(int64), args.Error(2)
}

func (m *MockEntityRepository) GetByDocument(ctx context.Context, document string) (*domain.Entity, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entity), args.Error(1)
}
