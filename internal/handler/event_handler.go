package handler

import (
	"net/http"
	"strconv"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/service"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EventHandler gerencia requisições de eventos
type EventHandler struct {
	service *service.EventService
	logger  *zap.Logger
}

// NewEventHandler cria um novo handler de eventos
func NewEventHandler(service *service.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger,
	}
}

// Create cria um novo evento
// POST /api/v1/events
func (h *EventHandler) Create(c *gin.Context) {
	// Obter organization_id do contexto (setado pelo middleware de auth)
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	// Obter user_id do contexto
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "user_id not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid user_id")
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	event, err := h.service.Create(c.Request.Context(), orgID, userID, &req)
	if err != nil {
		h.logger.Error("Failed to create event",
			zap.String("organization_id", orgIDStr.(string)),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to create event")
		return
	}

	response.Created(c, event)
}

// GetByID busca um evento por ID
// GET /api/v1/events/:id
func (h *EventHandler) GetByID(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	// Verificar se quer incluir participants
	includeParticipants := c.Query("include_participants") == "true"

	var event *dto.EventResponse
	if includeParticipants {
		event, err = h.service.GetByIDWithParticipants(c.Request.Context(), orgID, eventID)
	} else {
		event, err = h.service.GetByID(c.Request.Context(), orgID, eventID)
	}

	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to get event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to get event")
		return
	}

	response.Success(c, event)
}

// Update atualiza um evento
// PUT /api/v1/events/:id
func (h *EventHandler) Update(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	event, err := h.service.Update(c.Request.Context(), orgID, eventID, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to update event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to update event")
		return
	}

	response.Success(c, event)
}

// Delete remove um evento
// DELETE /api/v1/events/:id
func (h *EventHandler) Delete(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), orgID, eventID); err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to delete event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to delete event")
		return
	}

	response.NoContent(c)
}

// List lista eventos
// GET /api/v1/events
func (h *EventHandler) List(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	// Paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Filtro por status
	statusStr := c.Query("status")

	var events []*dto.EventResponse
	var total int64

	if statusStr != "" {
		status := domain.EventStatus(statusStr)
		events, total, err = h.service.ListByStatus(c.Request.Context(), orgID, status, page, perPage)
	} else {
		events, total, err = h.service.List(c.Request.Context(), orgID, page, perPage)
	}

	if err != nil {
		h.logger.Error("Failed to list events",
			zap.String("organization_id", orgIDStr.(string)),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to list events")
		return
	}

	response.Paginated(c, events, page, perPage, total)
}

// Activate ativa um evento
// POST /api/v1/events/:id/activate
func (h *EventHandler) Activate(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	event, err := h.service.Activate(c.Request.Context(), orgID, eventID)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to activate event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to activate event")
		return
	}

	response.Success(c, event)
}

// Cancel cancela um evento
// POST /api/v1/events/:id/cancel
func (h *EventHandler) Cancel(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	event, err := h.service.Cancel(c.Request.Context(), orgID, eventID)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to cancel event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to cancel event")
		return
	}

	response.Success(c, event)
}

// Complete marca um evento como completo
// POST /api/v1/events/:id/complete
func (h *EventHandler) Complete(c *gin.Context) {
	orgIDStr, exists := c.Get("organization_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "organization_id not found in context")
		return
	}

	orgID, err := uuid.Parse(orgIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid organization_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	event, err := h.service.Complete(c.Request.Context(), orgID, eventID)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		h.logger.Error("Failed to complete event",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to complete event")
		return
	}

	response.Success(c, event)
}
