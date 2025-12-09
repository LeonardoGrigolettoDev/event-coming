package handler

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
	startTime   time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		startTime:   time.Now(),
	}
}

// HealthStatus represents overall health status
type HealthStatus struct {
	Status    string                     `json:"status"`
	Service   string                     `json:"service"`
	Version   string                     `json:"version,omitempty"`
	Uptime    string                     `json:"uptime"`
	Timestamp time.Time                  `json:"timestamp"`
	Checks    map[string]ComponentHealth `json:"checks"`
}

// ComponentHealth represents health of a single component
type ComponentHealth struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	LatencyMs int64  `json:"latency_ms,omitempty"`
}

// Health returns basic health status (liveness probe)
// GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "event-coming",
	})
}

// Ready returns detailed health status (readiness probe)
// GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]ComponentHealth)
	allHealthy := true

	// Check PostgreSQL
	dbHealth := h.checkDatabase(ctx)
	checks["database"] = dbHealth
	if dbHealth.Status != "healthy" {
		allHealthy = false
	}

	// Check Redis
	redisHealth := h.checkRedis(ctx)
	checks["redis"] = redisHealth
	if redisHealth.Status != "healthy" {
		allHealthy = false
	}

	// Build response
	status := "healthy"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	response := HealthStatus{
		Status:    status,
		Service:   "event-coming",
		Version:   "1.0.0",
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now(),
		Checks:    checks,
	}

	c.JSON(httpStatus, response)
}

// Metrics returns basic metrics (could be expanded for Prometheus)
// GET /metrics
func (h *HealthHandler) Metrics(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(http.StatusOK, gin.H{
		"goroutines":     runtime.NumGoroutine(),
		"alloc_mb":       memStats.Alloc / 1024 / 1024,
		"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
		"sys_mb":         memStats.Sys / 1024 / 1024,
		"gc_cycles":      memStats.NumGC,
		"uptime_seconds": time.Since(h.startTime).Seconds(),
	})
}

// checkDatabase checks PostgreSQL connectivity
func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentHealth {
	if h.db == nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Message: "database not configured",
		}
	}

	start := time.Now()
	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	return ComponentHealth{
		Status:    "healthy",
		LatencyMs: time.Since(start).Milliseconds(),
	}
}

// checkRedis checks Redis connectivity
func (h *HealthHandler) checkRedis(ctx context.Context) ComponentHealth {
	if h.redisClient == nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Message: "redis not configured",
		}
	}

	start := time.Now()
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	return ComponentHealth{
		Status:    "healthy",
		LatencyMs: time.Since(start).Milliseconds(),
	}
}
