package router

import (
	"event-coming/internal/config"
	"event-coming/internal/handler"
	"event-coming/internal/handler/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router holds all dependencies needed for routing
type Router struct {
	engine             *gin.Engine
	config             *config.Config
	logger             *zap.Logger
	authHandler        *handler.AuthHandler
	websocketHandler   *handler.WebSocketHandler
	eventCacheHandler  *handler.EventCacheHandler
	participantHandler *handler.ParticipantHandler
	eventHandler       *handler.EventHandler
	entityHandler      *handler.EntityHandler
	locationHandler    *handler.LocationHandler
	webhookHandler     *handler.WebhookHandler
}

// NewRouter creates a new router
func NewRouter(
	cfg *config.Config,
	logger *zap.Logger,
	authHandler *handler.AuthHandler,
	websocketHandler *handler.WebSocketHandler,
	eventCacheHandler *handler.EventCacheHandler,
	participantHandler *handler.ParticipantHandler,
	eventHandler *handler.EventHandler,
	entityHandler *handler.EntityHandler,
	locationHandler *handler.LocationHandler,
	webhookHandler *handler.WebhookHandler,
) *Router {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	return &Router{
		engine:             engine,
		config:             cfg,
		logger:             logger,
		authHandler:        authHandler,
		websocketHandler:   websocketHandler,
		eventCacheHandler:  eventCacheHandler,
		participantHandler: participantHandler,
		eventHandler:       eventHandler,
		entityHandler:      entityHandler,
		locationHandler:    locationHandler,
		webhookHandler:     webhookHandler,
	}
}

// Setup configures all routes
func (r *Router) Setup() *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.RequestID())
	r.engine.Use(middleware.Recovery(r.logger))
	r.engine.Use(middleware.Logger(r.logger))
	r.engine.Use(middleware.CORS())

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "event-coming",
		})
	})

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.Refresh)
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/forgot-password", r.authHandler.ForgotPassword)
			auth.POST("/reset-password", r.authHandler.ResetPassword)
		}

		// WhatsApp webhook (public - called by WhatsApp servers)
		webhook := v1.Group("/webhook")
		{
			webhook.GET("/whatsapp", r.webhookHandler.VerifyWebhook)
			webhook.POST("/whatsapp", r.webhookHandler.HandleWebhook)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(&r.config.JWT))
		{
			// Entities
			entities := protected.Group("/entities")
			{
				entities.POST("", r.entityHandler.Create)
				entities.GET("", r.entityHandler.List)
				entities.GET("/:id", r.entityHandler.GetByID)
				entities.PUT("/:id", r.entityHandler.Update)
				entities.DELETE("/:id", r.entityHandler.Delete)
				entities.GET("/:id/children", r.entityHandler.ListByParent)
				entities.GET("/document/:document", r.entityHandler.GetByDocument)
			}

			// Events
			events := protected.Group("/events")
			{
				events.POST("", r.eventHandler.Create)
				events.GET("/:id", r.eventHandler.GetByID)
				events.PUT("/:id", r.eventHandler.Update)
				events.DELETE("/:id", r.eventHandler.Delete)
				events.GET("", r.eventHandler.List)

				// Event actions
				events.POST("/:id/activate", r.eventHandler.Activate)
				events.POST("/:id/cancel", r.eventHandler.Cancel)
				events.POST("/:id/complete", r.eventHandler.Complete)

				// Participants dentro de Events (usando :id consistente)
				events.POST("/:id/participants", r.participantHandler.Create)
				events.GET("/:id/participants", r.participantHandler.ListByEvent)
				events.POST("/:id/participants/batch", r.participantHandler.BatchCreate)

				// Locations for event (all participants)
				events.GET("/:id/locations", r.locationHandler.GetEventLocations)
			}

			// Participants
			participants := protected.Group("/participants")
			{
				participants.GET("/:id", r.participantHandler.GetByID)
				participants.PUT("/:id", r.participantHandler.Update)
				participants.DELETE("/:id", r.participantHandler.Delete)
				participants.POST("/:id/confirm", r.participantHandler.Confirm)
				participants.POST("/:id/check-in", r.participantHandler.CheckIn)

				// Locations
				participants.POST("/:id/locations", r.locationHandler.CreateLocation)
				participants.GET("/:id/locations", r.locationHandler.GetLocationHistory)
				participants.GET("/:id/locations/latest", r.locationHandler.GetLatestLocation)
			}

			// ETA
			eta := protected.Group("/eta")
			{
				eta.GET("/events/:id", r.locationHandler.GetEventETAs)
				eta.GET("/participants/:id", r.locationHandler.GetParticipantETA)
			}

			// Event cache (locations and confirmations from Redis) - movido para evitar conflito
			cache := protected.Group("/cache/:event")
			{
				cache.GET("", r.eventCacheHandler.GetEventCache)
				cache.GET("/locations", r.eventCacheHandler.GetLocationsOnly)
				cache.GET("/confirmations", r.eventCacheHandler.GetConfirmationsOnly)
			}
		}

		// WebSocket endpoint (fora do protected, autenticação via query param)
		v1.GET("/ws/:event", r.websocketHandler.HandleConnection)
	}

	return r.engine
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
