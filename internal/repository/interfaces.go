package repository

import (
	"context"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// OrganizationRepository defines organization data access methods
type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	Update(ctx context.Context, id uuid.UUID, input *domain.UpdateOrganizationInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, perPage int) ([]*domain.Organization, int64, error)
}

// UserRepository defines user data access methods
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User-Organization methods
	AddToOrganization(ctx context.Context, userOrg *domain.UserOrganization) error
	RemoveFromOrganization(ctx context.Context, userID, orgID uuid.UUID) error
	GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.UserOrganization, error)
	GetOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.User, error)
}

// EventRepository defines event data access methods
type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, orgID uuid.UUID, input *domain.UpdateEventInput) error
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*domain.Event, int64, error)
	ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.EventStatus, page, perPage int) ([]*domain.Event, int64, error)

	// Event instance methods
	CreateInstance(ctx context.Context, instance *domain.EventInstance) error
	GetInstanceByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.EventInstance, error)
	ListInstances(ctx context.Context, eventID uuid.UUID, orgID uuid.UUID) ([]*domain.EventInstance, error)
}

// ParticipantRepository defines participant data access methods
type ParticipantRepository interface {
	Create(ctx context.Context, participant *domain.Participant) error
	GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Participant, error)
	Update(ctx context.Context, id uuid.UUID, orgID uuid.UUID, input *domain.UpdateParticipantInput) error
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, orgID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error)
	ListByEventInstance(ctx context.Context, instanceID uuid.UUID, orgID uuid.UUID, page, perPage int) ([]*domain.Participant, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, orgID uuid.UUID, status domain.ParticipantStatus) error
	GetByPhoneNumber(ctx context.Context, phoneNumber string, eventID uuid.UUID, orgID uuid.UUID) (*domain.Participant, error)
}

// LocationRepository defines location data access methods
type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	BatchCreate(ctx context.Context, locations []*domain.Location) error
	GetLatestByParticipant(ctx context.Context, participantID uuid.UUID, orgID uuid.UUID) (*domain.Location, error)
	GetLatestByEvent(ctx context.Context, eventID uuid.UUID, orgID uuid.UUID) ([]*domain.Location, error)
	GetHistory(ctx context.Context, participantID uuid.UUID, orgID uuid.UUID, from, to time.Time) ([]*domain.Location, error)
}

// SchedulerRepository defines scheduler data access methods
type SchedulerRepository interface {
	Create(ctx context.Context, scheduler *domain.Scheduler) error
	GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Scheduler, error)
	Update(ctx context.Context, scheduler *domain.Scheduler) error
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	ListPending(ctx context.Context, before time.Time, limit int) ([]*domain.Scheduler, error)
	MarkAsProcessed(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID, orgID uuid.UUID, errorMsg string) error
	IncrementRetries(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
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
