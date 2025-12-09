package handler

import (
	"net/http"
	"strconv"

	"event-coming/internal/dto"
	"event-coming/internal/service"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ParticipantHandler gerencia requisições de participantes
type ParticipantHandler struct {
	service *service.ParticipantService
	logger  *zap.Logger
}

// NewParticipantHandler cria um novo handler de participantes
func NewParticipantHandler(service *service.ParticipantService, logger *zap.Logger) *ParticipantHandler {
	return &ParticipantHandler{
		service: service,
		logger:  logger,
	}
}

// Create cria um novo participante vinculado a um evento
// POST /api/v1/events/:event_id/participants
func (h *ParticipantHandler) Create(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	var req dto.CreateParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	participant, err := h.service.Create(c.Request.Context(), entityID, eventID, &req)
	if err != nil {
		h.logger.Error("Failed to create participant",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)

		if err.Error() == "event not found" {
			response.Error(c, http.StatusNotFound, "not_found", "event not found")
			return
		}
		if err.Error() == "participant with this phone number already exists in this event" {
			response.Error(c, http.StatusConflict, "conflict", err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to create participant")
		return
	}

	response.Created(c, participant)
}

// GetByID busca um participante por ID
// GET /api/v1/participants/:id
func (h *ParticipantHandler) GetByID(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	participantIDStr := c.Param("id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid participant_id")
		return
	}

	participant, err := h.service.GetByID(c.Request.Context(), entityID, participantID)
	if err != nil {
		h.logger.Error("Failed to get participant",
			zap.String("participant_id", participantIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusNotFound, "not_found", "participant not found")
		return
	}

	response.Success(c, participant)
}

// Update atualiza um participante
// PUT /api/v1/participants/:id
func (h *ParticipantHandler) Update(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	participantIDStr := c.Param("id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid participant_id")
		return
	}

	var req dto.UpdateParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	participant, err := h.service.Update(c.Request.Context(), entityID, participantID, &req)
	if err != nil {
		h.logger.Error("Failed to update participant",
			zap.String("participant_id", participantIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to update participant")
		return
	}

	response.Success(c, participant)
}

// Delete remove um participante
// DELETE /api/v1/participants/:id
func (h *ParticipantHandler) Delete(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	participantIDStr := c.Param("id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid participant_id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), entityID, participantID); err != nil {
		h.logger.Error("Failed to delete participant",
			zap.String("participant_id", participantIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to delete participant")
		return
	}

	response.NoContent(c)
}

// ListByEvent lista participantes de um evento
// GET /api/v1/events/:event_id/participants
func (h *ParticipantHandler) ListByEvent(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
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

	participants, total, err := h.service.ListByEvent(c.Request.Context(), entityID, eventID, page, perPage)
	if err != nil {
		h.logger.Error("Failed to list participants",
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to list participants")
		return
	}

	response.Paginated(c, participants, page, perPage, total)
}

// Confirm confirma a participação
// POST /api/v1/participants/:id/confirm
func (h *ParticipantHandler) Confirm(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	participantIDStr := c.Param("id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid participant_id")
		return
	}

	participant, err := h.service.ConfirmParticipant(c.Request.Context(), entityID, participantID)
	if err != nil {
		h.logger.Error("Failed to confirm participant",
			zap.String("participant_id", participantIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to confirm participant")
		return
	}

	response.Success(c, participant)
}

// CheckIn faz check-in do participante
// POST /api/v1/participants/:id/check-in
func (h *ParticipantHandler) CheckIn(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	participantIDStr := c.Param("id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid participant_id")
		return
	}

	participant, err := h.service.CheckInParticipant(c.Request.Context(), entityID, participantID)
	if err != nil {
		h.logger.Error("Failed to check-in participant",
			zap.String("participant_id", participantIDStr),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "internal_error", "failed to check-in participant")
		return
	}

	response.Success(c, participant)
}

// BatchCreate cria múltiplos participantes
// POST /api/v1/events/:event_id/participants/batch
func (h *ParticipantHandler) BatchCreate(c *gin.Context) {
	entityIDStr, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusBadRequest, "bad_request", "entity_id not found in context")
		return
	}

	entityID, err := uuid.Parse(entityIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid entity_id")
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid event_id")
		return
	}

	var req dto.BatchCreateParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	participants, errors := h.service.BatchCreate(c.Request.Context(), entityID, eventID, &req)

	// Preparar resposta
	errorMessages := make([]string, len(errors))
	for i, err := range errors {
		errorMessages[i] = err.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"created":      len(participants),
		"failed":       len(errors),
		"participants": participants,
		"errors":       errorMessages,
	})
}
