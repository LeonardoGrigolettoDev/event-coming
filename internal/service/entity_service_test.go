package service

import (
	"context"
	"testing"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/testutil/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEntityService_Create(t *testing.T) {
	tests := []struct {
		name      string
		req       *dto.CreateEntityRequest
		setupMock func(*mocks.MockEntityRepository)
		wantErr   error
	}{
		{
			name: "successful creation",
			req: &dto.CreateEntityRequest{
				Type:  domain.EntityTypeNaturalPerson,
				Name:  "Test Entity",
				Email: entityStrPtr("test@example.com"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Entity")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "with document check - no conflict",
			req: &dto.CreateEntityRequest{
				Type:     domain.EntityTypeNaturalPerson,
				Name:     "Test Entity",
				Document: entityStrPtr("12345678901"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("GetByDocument", mock.Anything, "12345678901").Return(nil, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Entity")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "document already exists",
			req: &dto.CreateEntityRequest{
				Type:     domain.EntityTypeNaturalPerson,
				Name:     "Test Entity",
				Document: entityStrPtr("12345678901"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				existing := &domain.Entity{ID: uuid.New()}
				m.On("GetByDocument", mock.Anything, "12345678901").Return(existing, nil)
			},
			wantErr: domain.ErrConflict,
		},
		{
			name: "with parent - valid parent",
			req: &dto.CreateEntityRequest{
				Type:     domain.EntityTypeNaturalPerson,
				Name:     "Child Entity",
				ParentID: entityUUIDPtr(uuid.New()),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				parent := &domain.Entity{ID: uuid.New()}
				m.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(parent, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Entity")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "with parent - parent not found",
			req: &dto.CreateEntityRequest{
				Type:     domain.EntityTypeNaturalPerson,
				Name:     "Child Entity",
				ParentID: entityUUIDPtr(uuid.New()),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo)

			svc := NewEntityService(mockRepo)
			result, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*mocks.MockEntityRepository, uuid.UUID)
		wantErr   error
	}{
		{
			name: "entity found",
			id:   uuid.New(),
			setupMock: func(m *mocks.MockEntityRepository, id uuid.UUID) {
				entity := &domain.Entity{
					ID:   id,
					Name: "Test Entity",
					Type: domain.EntityTypeNaturalPerson,
				}
				m.On("GetByID", mock.Anything, id).Return(entity, nil)
			},
			wantErr: nil,
		},
		{
			name: "entity not found",
			id:   uuid.New(),
			setupMock: func(m *mocks.MockEntityRepository, id uuid.UUID) {
				m.On("GetByID", mock.Anything, id).Return(nil, nil)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo, tt.id)

			svc := NewEntityService(mockRepo)
			result, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_Update(t *testing.T) {
	entityID := uuid.New()
	parentID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		req       *dto.UpdateEntityRequest
		setupMock func(*mocks.MockEntityRepository)
		wantErr   error
	}{
		{
			name: "successful update",
			id:   entityID,
			req: &dto.UpdateEntityRequest{
				Name: entityStrPtr("Updated Name"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				entity := &domain.Entity{ID: entityID, Name: "Original Name"}
				m.On("GetByID", mock.Anything, entityID).Return(entity, nil).Once()
				m.On("Update", mock.Anything, entityID, mock.AnythingOfType("*domain.UpdateEntityInput")).Return(nil)
				m.On("GetByID", mock.Anything, entityID).Return(&domain.Entity{ID: entityID, Name: "Updated Name"}, nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "entity not found",
			id:   entityID,
			req: &dto.UpdateEntityRequest{
				Name: entityStrPtr("Updated Name"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("GetByID", mock.Anything, entityID).Return(nil, nil)
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "document conflict",
			id:   entityID,
			req: &dto.UpdateEntityRequest{
				Document: entityStrPtr("12345678901"),
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				entity := &domain.Entity{ID: entityID}
				otherEntity := &domain.Entity{ID: uuid.New()}
				m.On("GetByID", mock.Anything, entityID).Return(entity, nil)
				m.On("GetByDocument", mock.Anything, "12345678901").Return(otherEntity, nil)
			},
			wantErr: domain.ErrConflict,
		},
		{
			name: "self reference parent",
			id:   entityID,
			req: &dto.UpdateEntityRequest{
				ParentID: &entityID,
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				entity := &domain.Entity{ID: entityID}
				m.On("GetByID", mock.Anything, entityID).Return(entity, nil)
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "parent not found",
			id:   entityID,
			req: &dto.UpdateEntityRequest{
				ParentID: &parentID,
			},
			setupMock: func(m *mocks.MockEntityRepository) {
				entity := &domain.Entity{ID: entityID}
				m.On("GetByID", mock.Anything, entityID).Return(entity, nil).Once()
				m.On("GetByID", mock.Anything, parentID).Return(nil, nil).Once()
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo)

			svc := NewEntityService(mockRepo)
			result, err := svc.Update(context.Background(), tt.id, tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*mocks.MockEntityRepository, uuid.UUID)
		wantErr   error
	}{
		{
			name: "successful delete",
			id:   uuid.New(),
			setupMock: func(m *mocks.MockEntityRepository, id uuid.UUID) {
				entity := &domain.Entity{ID: id}
				m.On("GetByID", mock.Anything, id).Return(entity, nil)
				m.On("Delete", mock.Anything, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "entity not found",
			id:   uuid.New(),
			setupMock: func(m *mocks.MockEntityRepository, id uuid.UUID) {
				m.On("GetByID", mock.Anything, id).Return(nil, nil)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo, tt.id)

			svc := NewEntityService(mockRepo)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_List(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		perPage    int
		setupMock  func(*mocks.MockEntityRepository)
		wantLen    int
		wantTotal  int64
		wantErr    bool
		expectedPP int // expected perPage after normalization
	}{
		{
			name:    "successful list",
			page:    1,
			perPage: 10,
			setupMock: func(m *mocks.MockEntityRepository) {
				entities := []*domain.Entity{
					{ID: uuid.New(), Name: "Entity 1"},
					{ID: uuid.New(), Name: "Entity 2"},
				}
				m.On("List", mock.Anything, 1, 10).Return(entities, int64(2), nil)
			},
			wantLen:   2,
			wantTotal: 2,
		},
		{
			name:    "pagination normalization - page 0",
			page:    0,
			perPage: 10,
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("List", mock.Anything, 1, 10).Return([]*domain.Entity{}, int64(0), nil)
			},
			wantLen:   0,
			wantTotal: 0,
		},
		{
			name:    "pagination normalization - perPage 0",
			page:    1,
			perPage: 0,
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("List", mock.Anything, 1, 20).Return([]*domain.Entity{}, int64(0), nil)
			},
			wantLen:   0,
			wantTotal: 0,
		},
		{
			name:    "pagination normalization - perPage > 100",
			page:    1,
			perPage: 200,
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("List", mock.Anything, 1, 20).Return([]*domain.Entity{}, int64(0), nil)
			},
			wantLen:   0,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo)

			svc := NewEntityService(mockRepo)
			result, total, err := svc.List(context.Background(), tt.page, tt.perPage)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				assert.Equal(t, tt.wantTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_ListByParent(t *testing.T) {
	parentID := uuid.New()

	tests := []struct {
		name      string
		parentID  uuid.UUID
		page      int
		perPage   int
		setupMock func(*mocks.MockEntityRepository)
		wantLen   int
		wantTotal int64
	}{
		{
			name:     "successful list by parent",
			parentID: parentID,
			page:     1,
			perPage:  10,
			setupMock: func(m *mocks.MockEntityRepository) {
				entities := []*domain.Entity{
					{ID: uuid.New(), ParentID: &parentID, Name: "Child 1"},
				}
				m.On("ListByParent", mock.Anything, parentID, 1, 10).Return(entities, int64(1), nil)
			},
			wantLen:   1,
			wantTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo)

			svc := NewEntityService(mockRepo)
			result, total, err := svc.ListByParent(context.Background(), tt.parentID, tt.page, tt.perPage)

			assert.NoError(t, err)
			assert.Len(t, result, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEntityService_GetByDocument(t *testing.T) {
	tests := []struct {
		name      string
		document  string
		setupMock func(*mocks.MockEntityRepository)
		wantErr   error
	}{
		{
			name:     "entity found by document",
			document: "12345678901",
			setupMock: func(m *mocks.MockEntityRepository) {
				entity := &domain.Entity{
					ID:       uuid.New(),
					Name:     "Test Entity",
					Document: entityStrPtr("12345678901"),
				}
				m.On("GetByDocument", mock.Anything, "12345678901").Return(entity, nil)
			},
			wantErr: nil,
		},
		{
			name:     "entity not found by document",
			document: "12345678901",
			setupMock: func(m *mocks.MockEntityRepository) {
				m.On("GetByDocument", mock.Anything, "12345678901").Return(nil, nil)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockEntityRepository)
			tt.setupMock(mockRepo)

			svc := NewEntityService(mockRepo)
			result, err := svc.GetByDocument(context.Background(), tt.document)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions for entity service tests
func entityStrPtr(s string) *string {
	return &s
}

func entityUUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
