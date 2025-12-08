package handler

import (
	"net/http"

	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EventCacheHandler gerencia requisições de cache de eventos
type EventCacheHandler struct {
	service *service.EventCacheService
	logger  *zap.Logger
}

// NewEventCacheHandler cria um novo handler de cache de eventos
func NewEventCacheHandler(service *service.EventCacheService, logger *zap.Logger) *EventCacheHandler {
	return &EventCacheHandler{
		service: service,
		logger:  logger,
	}
}

// GetEventCache busca informações de localização e confirmações do cache
// GET /api/v1/:organization/:event
func (h *EventCacheHandler) GetEventCache(c *gin.Context) {
	orgIDStr := c.Param("organization")
	eventIDStr := c.Param("event")

	// Validar UUIDs
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid organization_id",
			"message": "organization_id must be a valid UUID",
		})
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid event_id",
			"message": "event_id must be a valid UUID",
		})
		return
	}

	// Buscar dados do cache
	data, err := h.service.GetEventCacheData(c.Request.Context(), orgID, eventID)
	if err != nil {
		h.logger.Error("Failed to get event cache data",
			zap.String("organization_id", orgIDStr),
			zap.String("event_id", eventIDStr),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve cache data",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetLocationsOnly busca apenas as localizações do cache
// GET /api/v1/:organization/:event/locations
func (h *EventCacheHandler) GetLocationsOnly(c *gin.Context) {
	orgIDStr := c.Param("organization")
	eventIDStr := c.Param("event")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization_id"})
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	data, err := h.service.GetEventCacheData(c.Request.Context(), orgID, eventID)
	if err != nil {
		h.logger.Error("Failed to get locations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id": orgID,
		"event_id":        eventID,
		"locations":       data.Locations,
		"total":           data.TotalLocations,
	})
}

// GetConfirmationsOnly busca apenas as confirmações do cache
// GET /api/v1/:organization/:event/confirmations
func (h *EventCacheHandler) GetConfirmationsOnly(c *gin.Context) {
	orgIDStr := c.Param("organization")
	eventIDStr := c.Param("event")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization_id"})
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	data, err := h.service.GetEventCacheData(c.Request.Context(), orgID, eventID)
	if err != nil {
		h.logger.Error("Failed to get confirmations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id": orgID,
		"event_id":        eventID,
		"confirmations":   data.Confirmations,
		"total_confirmed": data.TotalConfirmed,
		"total_pending":   data.TotalPending,
		"total_denied":    data.TotalDenied,
	})
}
