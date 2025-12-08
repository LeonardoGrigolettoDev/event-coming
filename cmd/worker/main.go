package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-coming/internal/cache"
	"event-coming/internal/config"
	"event-coming/internal/repository/postgres"
	"event-coming/internal/service"
	"event-coming/internal/whatsapp"
	"event-coming/internal/worker"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Event-Coming Workers")

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

	// Connect to Redis
	logger.Info("Connecting to Redis")
	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()
	logger.Info("Connected to Redis")

	// Initialize repositories
	schedulerRepo := postgres.NewSchedulerRepository(db)
	participantRepo := postgres.NewParticipantRepository(db)
	eventRepo := postgres.NewEventRepository(db)

	// Initialize WhatsApp client (pode ser nil se não configurado)
	var whatsappClient *whatsapp.Client
	if cfg.WhatsApp.AccessToken != "" {
		whatsappClient = whatsapp.NewClient(&cfg.WhatsApp)
		logger.Info("WhatsApp client initialized")
	} else {
		logger.Warn("WhatsApp client not configured, notifications will be skipped")
	}

	// Initialize services
	notificationService := service.NewNotificationService(whatsappClient, logger)
	schedulerService := service.NewSchedulerService(
		schedulerRepo,
		participantRepo,
		eventRepo,
		notificationService,
		logger,
	)

	// Initialize workers
	schedulerWorker := worker.NewSchedulerWorker(
		schedulerService,
		logger,
		30*time.Second, // Intervalo de processamento
		100,            // Batch size
	)

	// Start workers in goroutines
	go schedulerWorker.Start(ctx)

	logger.Info("All workers started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down workers...")

	// Cancel context to signal workers to stop
	cancel()

	// Stop workers gracefully
	schedulerWorker.Stop()

	logger.Info("Workers exited gracefully")
}
