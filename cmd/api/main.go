package main

import (
	"context"
	"event-coming/internal/cache"
	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/internal/handler"
	"event-coming/internal/repository/postgres"
	"event-coming/internal/router"
	"event-coming/internal/service"
	"event-coming/internal/websocket"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Event-Coming API")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to PostgreSQL
	logger.Info("Connecting to PostgreSQL")
	db, err := postgres.NewGormDB(&cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	logger.Info("Connected to PostgreSQL")

	if cfg.App.Debug {
		logger.Info("Running AutoMigrate (dev mode)...")
		db.AutoMigrate(
			&domain.User{},
			&domain.RefreshToken{},
			&domain.Entity{},
			&domain.Entity{},
			&domain.Participant{},
			&domain.Event{},
			&domain.EventInstance{},
			&domain.UserEntity{},
			&domain.Location{},
			&domain.Scheduler{},
		)
	}

	// Connect to Redis
	logger.Info("Connecting to Redis")
	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()
	logger.Info("Connected to Redis")

	// Initialize WebSocket Hub and PubSub
	wsHub := websocket.NewHub(logger)
	wsPubSub := websocket.NewPubSub(redisClient, wsHub, logger)

	// Start WebSocket Hub
	go wsHub.Run(ctx)

	// Subscribe to all Redis channels
	if err := wsPubSub.SubscribeAll(ctx); err != nil {
		logger.Warn("Failed to subscribe to Redis PubSub", zap.Error(err))
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewRefreshTokenRepository(db)
	participantRepo := postgres.NewParticipantRepository(db)
	eventRepo := postgres.NewEventRepository(db)
	schedulerRepo := postgres.NewSchedulerRepository(db)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		&cfg.JWT,
	)
	eventCacheService := service.NewEventCacheService(redisClient)
	participantService := service.NewParticipantService(participantRepo, eventRepo)
	eventService := service.NewEventService(eventRepo, schedulerRepo, participantRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	websocketHandler := handler.NewWebSocketHandler(wsHub, wsPubSub, logger)
	eventCacheHandler := handler.NewEventCacheHandler(eventCacheService, logger)
	participantHandler := handler.NewParticipantHandler(participantService, logger)
	eventHandler := handler.NewEventHandler(eventService, logger)

	// Setup router
	r := router.NewRouter(cfg, logger, authHandler, websocketHandler, eventCacheHandler, participantHandler, eventHandler)
	engine := r.Setup()

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
