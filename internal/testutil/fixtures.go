package testutil

import (
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
)

// Test UUIDs
var (
	TestUserID        = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	TestEntityID      = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	TestEventID       = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	TestParticipantID = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	TestSchedulerID   = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	TestLocationID    = uuid.MustParse("66666666-6666-6666-6666-666666666666")
)

// NewTestUser creates a test user
func NewTestUser() *domain.User {
	phone := "+5511999999999"
	return &domain.User{
		ID:            TestUserID,
		Email:         "test@example.com",
		PasswordHash:  "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqJXZu1Z1Z1Z1Z1Z1Z1Z1Z1Z1Z1Z1", // hashed "password123"
		Name:          "Test User",
		Phone:         &phone,
		Active:        true,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewTestEntity creates a test entity
func NewTestEntity() *domain.Entity {
	email := "entity@example.com"
	phone := "+5511888888888"
	doc := "12345678901"
	return &domain.Entity{
		ID:          TestEntityID,
		Type:        domain.EntityTypeLegalEntity,
		Name:        "Test Company",
		Email:       &email,
		PhoneNumber: &phone,
		Document:    &doc,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewTestEvent creates a test event
func NewTestEvent() *domain.Event {
	lat := -23.550520
	lng := -46.633308
	addr := "São Paulo, SP"
	endTime := time.Now().Add(4 * time.Hour)
	desc := "A test event description"
	return &domain.Event{
		ID:              TestEventID,
		EntityID:        TestEntityID,
		Name:            "Test Event",
		Description:     &desc,
		Status:          domain.EventStatusDraft,
		LocationLat:     lat,
		LocationLng:     lng,
		LocationAddress: &addr,
		StartTime:       time.Now().Add(2 * time.Hour),
		EndTime:         &endTime,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// NewTestParticipant creates a test participant
func NewTestParticipant() *domain.Participant {
	return &domain.Participant{
		ID:        TestParticipantID,
		EventID:   TestEventID,
		EntityID:  TestEntityID,
		Status:    domain.ParticipantStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestScheduler creates a test scheduler
func NewTestScheduler() *domain.Scheduler {
	return &domain.Scheduler{
		ID:          TestSchedulerID,
		EntityID:    TestEntityID,
		EventID:     TestEventID,
		Action:      domain.SchedulerActionConfirmation,
		Status:      domain.SchedulerStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Hour),
		MaxRetries:  3,
		Retries:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewTestLocation creates a test location
func NewTestLocation() *domain.Location {
	acc := 10.5
	alt := 760.0
	speed := 25.5
	head := 180.0
	return &domain.Location{
		ID:            TestLocationID,
		ParticipantID: TestParticipantID,
		EventID:       TestEventID,
		EntityID:      TestEntityID,
		Latitude:      -23.561684,
		Longitude:     -46.655981,
		Accuracy:      &acc,
		Altitude:      &alt,
		Speed:         &speed,
		Heading:       &head,
		Timestamp:     time.Now(),
	}
}

// NewTestRefreshToken creates a test refresh token
func NewTestRefreshToken() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    TestUserID,
		Token:     "test_token_hash_12345",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

// NewTestPasswordResetToken creates a test password reset token
func NewTestPasswordResetToken() *domain.PasswordResetToken {
	return &domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    TestUserID,
		Token:     "reset_token_hash_12345",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}
}
