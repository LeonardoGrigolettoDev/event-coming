package repository

import (
	"context"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// EntityRepository defines entity data access methods
type EntityRepository interface {
	Create(ctx context.Context, entity *domain.Entity) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error)
	Update(ctx context.Context, id uuid.UUID, input *domain.UpdateEntityInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, perPage int) ([]*domain.Entity, int64, error)
	ListByParent(ctx context.Context, parentID uuid.UUID, page, perPage int) ([]*domain.Entity, int64, error)
	GetByDocument(ctx context.Context, document string) (*domain.Entity, error)
}

// UserRepository defines user data access methods
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User-Entity methods
	AddToEntity(ctx context.Context, userEntity *domain.UserEntity) error
	RemoveFromEntity(ctx context.Context, userID, entityID uuid.UUID) error
	GetUserEntities(ctx context.Context, userID uuid.UUID) ([]*domain.UserEntity, error)
	GetEntityUsers(ctx context.Context, entityID uuid.UUID) ([]*domain.User, error)
}

// EventRepository defines event data access methods
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateEventInput) error
	Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
	List(ctx context.Context, entityID uuid.UUID, page, perPage int) ([]*domain.Event, int64, error)
	ListByStatus(ctx context.Context, entityID uuid.UUID, status domain.EventStatus, page, perPage int) ([]*domain.Event, int64, error)

	// Event instance methods
	CreateInstance(ctx context.Context, instance *domain.EventInstance) error
	GetInstanceByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.EventInstance, error)
	ListInstances(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.EventInstance, error)
}

// ParticipantRepository defines participant data access methods
type ParticipantRepository interface {
	Create(ctx context.Context, participant *domain.Participant) error
	GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Participant, error)
	Update(ctx context.Context, id uuid.UUID, entityID uuid.UUID, input *domain.UpdateParticipantInput) error
	Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error)
	ListByEventInstance(ctx context.Context, instanceID uuid.UUID, entityID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, entityID uuid.UUID, status domain.ParticipantStatus) error
	GetByPhoneNumber(ctx context.Context, phoneNumber string, eventID uuid.UUID, entityID uuid.UUID) (*domain.Participant, error)
	// GetActiveByPhoneNumber finds a participant by phone number in active events
	GetActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Participant, error)
}

// LocationRepository defines location data access methods
type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	BatchCreate(ctx context.Context, locations []*domain.Location) error
	GetLatestByParticipant(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID) (*domain.Location, error)
	GetLatestByEvent(ctx context.Context, eventID uuid.UUID, entityID uuid.UUID) ([]*domain.Location, error)
	GetHistory(ctx context.Context, participantID uuid.UUID, entityID uuid.UUID, from, to time.Time) ([]*domain.Location, error)
}

// SchedulerRepository defines scheduler data access methods
type SchedulerRepository interface {
	Create(ctx context.Context, scheduler *domain.Scheduler) error
	GetByID(ctx context.Context, id uuid.UUID, entityID uuid.UUID) (*domain.Scheduler, error)
	Update(ctx context.Context, scheduler *domain.Scheduler) error
	Delete(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
	ListPending(ctx context.Context, before time.Time, limit int) ([]*domain.Scheduler, error)
	MarkAsProcessed(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID, entityID uuid.UUID, errorMsg string) error
	IncrementRetries(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
}

// RefreshTokenRepository defines refresh token data access methods
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)

	// Revogar por ID (interno, após refresh)
	Revoke(ctx context.Context, id uuid.UUID) error

	// Revogar por hash do token (logout)
	RevokeByToken(ctx context.Context, tokenHash string) error

	// Revogar todos do usuário (reset password, segurança)
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error

	DeleteExpired(ctx context.Context) error
}

// PasswordResetTokenRepository defines password reset token data access methods
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *domain.PasswordResetToken) error
	GetByToken(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
