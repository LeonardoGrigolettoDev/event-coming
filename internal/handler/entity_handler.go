package handler

import (
	"net/http"
	"strconv"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/service"
	"event-coming/pkg/response"
	"event-coming/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EntityHandler handles entity HTTP requests
type EntityHandler struct {
	entityService *service.EntityService
	logger        *zap.Logger
}

// NewEntityHandler creates a new entity handler
func NewEntityHandler(entityService *service.EntityService, logger *zap.Logger) *EntityHandler {
	return &EntityHandler{
		entityService: entityService,
		logger:        logger,
	}
}

// Create handles POST /entities
func (h *EntityHandler) Create(c *gin.Context) {
	var req dto.CreateEntityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	entity, err := h.entityService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create entity", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.Created(c, entity)
}

// GetByID handles GET /entities/:id
func (h *EntityHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid entity ID")
		return
	}

	entity, err := h.entityService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get entity", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.Success(c, entity)
}

// Update handles PUT /entities/:id
func (h *EntityHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid entity ID")
		return
	}

	var req dto.UpdateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	entity, err := h.entityService.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update entity", zap.Error(err))
		if err == domain.ErrInvalidInput {
			response.Error(c, http.StatusBadRequest, "bad_request", "Entity cannot be its own parent")
			return
		}
		response.HandleDomainError(c, err)
		return
	}

	response.Success(c, entity)
}

// Delete handles DELETE /entities/:id
func (h *EntityHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid entity ID")
		return
	}

	if err := h.entityService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete entity", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.NoContent(c)
}

// List handles GET /entities
func (h *EntityHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	entities, total, err := h.entityService.List(c.Request.Context(), page, perPage)
	if err != nil {
		h.logger.Error("Failed to list entities", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.Paginated(c, entities, page, perPage, total)
}

// ListByParent handles GET /entities/:id/children
func (h *EntityHandler) ListByParent(c *gin.Context) {
	idStr := c.Param("id")
	parentID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid parent entity ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	entities, total, err := h.entityService.ListByParent(c.Request.Context(), parentID, page, perPage)
	if err != nil {
		h.logger.Error("Failed to list child entities", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.Paginated(c, entities, page, perPage, total)
}

// GetByDocument handles GET /entities/document/:document
func (h *EntityHandler) GetByDocument(c *gin.Context) {
	document := c.Param("document")
	if document == "" {
		response.Error(c, http.StatusBadRequest, "bad_request", "Document is required")
		return
	}

	entity, err := h.entityService.GetByDocument(c.Request.Context(), document)
	if err != nil {
		h.logger.Error("Failed to get entity by document", zap.Error(err))
		response.HandleDomainError(c, err)
		return
	}

	response.Success(c, entity)
}
