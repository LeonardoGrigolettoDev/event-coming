package handler

import (
	"net/http"
	"time"

	"event-coming/internal/domain"
	"event-coming/internal/dto"
	"event-coming/internal/service"
	"event-coming/internal/service/eta"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LocationHandler handles location-related HTTP requests
type LocationHandler struct {
	locationService *service.LocationService
	etaService      *eta.ETAService
	eventService    *service.EventService
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(
	locationService *service.LocationService,
	etaService *eta.ETAService,
	eventService *service.EventService,
) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
		etaService:      etaService,
		eventService:    eventService,
	}
}

// CreateLocation creates a new location for a participant
// POST /participants/:id/locations
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid participant ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.locationService.CreateLocation(c.Request.Context(), participantID, entityID.(uuid.UUID), &req)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "Participant not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response.Created(c, result)
}

// GetLocationHistory gets location history for a participant
// GET /participants/:id/locations
func (h *LocationHandler) GetLocationHistory(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid participant ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	// Parse time range (default: last 1 hour)
	now := time.Now()
	from := now.Add(-1 * time.Hour)
	to := now

	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = parsed
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = parsed
		}
	}

	locations, err := h.locationService.GetLocationHistory(
		c.Request.Context(),
		participantID,
		entityID.(uuid.UUID),
		from,
		to,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response.Success(c, locations)
}

// GetLatestLocation gets the latest location for a participant
// GET /participants/:id/locations/latest
func (h *LocationHandler) GetLatestLocation(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid participant ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	location, err := h.locationService.GetLatestLocation(
		c.Request.Context(),
		participantID,
		entityID.(uuid.UUID),
	)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "Location not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response.Success(c, location)
}

// GetEventLocations gets latest locations for all participants in an event
// GET /events/:id/locations
func (h *LocationHandler) GetEventLocations(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid event ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	locations, err := h.locationService.GetEventLocations(
		c.Request.Context(),
		eventID,
		entityID.(uuid.UUID),
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response.Success(c, locations)
}

// GetParticipantETA gets ETA for a participant to reach event location
// GET /eta/participants/:id
func (h *LocationHandler) GetParticipantETA(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid participant ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	eventIDStr := c.Query("event_id")
	if eventIDStr == "" {
		response.Error(c, http.StatusBadRequest, "bad_request", "event_id query parameter is required")
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid event ID")
		return
	}

	// Get event to get target location
	event, err := h.eventService.GetByID(c.Request.Context(), entityID.(uuid.UUID), eventID)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	if event.LocationLat == 0 && event.LocationLng == 0 {
		response.Error(c, http.StatusBadRequest, "bad_request", "Event does not have a location defined")
		return
	}

	result, err := h.etaService.CalculateETA(
		c.Request.Context(),
		participantID,
		entityID.(uuid.UUID),
		event.LocationLat,
		event.LocationLng,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	etaResponse := dto.ToETAResponse(result)
	response.Success(c, etaResponse)
}

// GetEventETAs gets ETAs for all participants in an event
// GET /eta/events/:id
func (h *LocationHandler) GetEventETAs(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", "Invalid event ID")
		return
	}

	entityID, exists := c.Get("entity_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "Entity not found in context")
		return
	}

	// Get event with participants
	event, err := h.eventService.GetByIDWithParticipants(c.Request.Context(), entityID.(uuid.UUID), eventID)
	if err != nil {
		if err == domain.ErrNotFound {
			response.Error(c, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	if event.LocationLat == 0 && event.LocationLng == 0 {
		response.Error(c, http.StatusBadRequest, "bad_request", "Event does not have a location defined")
		return
	}

	// Collect participant IDs
	participantIDs := make([]uuid.UUID, len(event.Participants))
	for i, p := range event.Participants {
		participantIDs[i] = p.ID
	}

	// Calculate ETAs for all participants
	results, err := h.etaService.CalculateMultipleETAs(
		c.Request.Context(),
		participantIDs,
		entityID.(uuid.UUID),
		event.LocationLat,
		event.LocationLng,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	eventETA := dto.ToEventETAResponse(eventID, entityID.(uuid.UUID), results)
	response.Success(c, eventETA)
}
