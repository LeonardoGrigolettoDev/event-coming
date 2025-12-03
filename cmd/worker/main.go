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

	// Create context with timeout for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to PostgreSQL
	logger.Info("Connecting to PostgreSQL")
	db, err := postgres.NewPool(ctx, &cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()
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
	// Note: Add repository initialization here
	// Example:
	// schedulerRepo := postgres.NewSchedulerRepository(db)
	// locationRepo := postgres.NewLocationRepository(db)

	// Initialize workers
	// Note: Add worker initialization and startup here
	// Example:
	// schedulerWorker := worker.NewSchedulerWorker(schedulerRepo, whatsappClient, logger)
	// locationFlusher := worker.NewLocationFlusher(locationBuffer, locationRepo, logger)
	// recurrenceWorker := worker.NewRecurrenceWorker(eventRepo, logger)

	// Start workers in goroutines
	// go schedulerWorker.Start(ctx)
	// go locationFlusher.Start(ctx)
	// go recurrenceWorker.Start(ctx)

	logger.Info("All workers started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down workers...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop workers
	// Note: Implement graceful shutdown for each worker
	// Example:
	// schedulerWorker.Stop(shutdownCtx)
	// locationFlusher.Stop(shutdownCtx)
	// recurrenceWorker.Stop(shutdownCtx)
	
	// Use context for potential cleanup
	_ = shutdownCtx

	logger.Info("Workers exited gracefully")
}
