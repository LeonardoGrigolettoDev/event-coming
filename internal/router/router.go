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
			auth.POST("/register", r.authHandler.Register) // ← Direto!
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.Refresh)
			auth.POST("/forgot-password", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "not implemented"})
			})
			auth.POST("/reset-password", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "not implemented"})
			})
		}

		// WhatsApp webhook
		webhook := v1.Group("/webhook")
		{
			webhook.GET("/whatsapp", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "not implemented"})
			})
			webhook.POST("/whatsapp", func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "not implemented"})
			})
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(&r.config.JWT))
		{
			// Organizations
			orgs := protected.Group("/organizations")
			{
				orgs.POST("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
				orgs.GET("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
				orgs.PUT("/:id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
				orgs.GET("", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
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
				participants.POST("/:id/locations", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
				participants.GET("/:id/locations", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
			}

			// ETA
			eta := protected.Group("/eta")
			{
				eta.GET("/events/:event_id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
				eta.GET("/participants/:participant_id", func(c *gin.Context) {
					c.JSON(501, gin.H{"message": "not implemented"})
				})
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
