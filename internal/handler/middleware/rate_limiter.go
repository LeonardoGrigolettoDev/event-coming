package middleware

import (
	"net/http"
	"sync"
	"time"

	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	// Requests per second allowed
	RequestsPerSecond float64
	// Burst size (max tokens)
	BurstSize int
	// Cleanup interval for expired entries
	CleanupInterval time.Duration
}

// DefaultRateLimiterConfig returns sensible defaults
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute * 5,
	}
}

// tokenBucket implements the token bucket algorithm
type tokenBucket struct {
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// take tries to take a token from the bucket
func (b *tokenBucket) take(rate float64, burst int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastUpdate).Seconds()
	b.lastUpdate = now

	// Add tokens based on elapsed time
	b.tokens += elapsed * rate
	if b.tokens > float64(burst) {
		b.tokens = float64(burst)
	}

	// Try to consume a token
	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

// RateLimiter stores rate limit state for each IP
type RateLimiter struct {
	buckets  map[string]*tokenBucket
	mu       sync.RWMutex
	config   RateLimiterConfig
	stopChan chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		config:   config,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given key is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(rl.config.BurstSize),
			lastUpdate: time.Now(),
		}
		rl.buckets[key] = bucket
	}
	rl.mu.Unlock()

	return bucket.take(rl.config.RequestsPerSecond, rl.config.BurstSize)
}

// cleanup removes expired buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, bucket := range rl.buckets {
				bucket.mu.Lock()
				// Remove if no activity for 2x cleanup interval
				if now.Sub(bucket.lastUpdate) > rl.config.CleanupInterval*2 {
					delete(rl.buckets, key)
				}
				bucket.mu.Unlock()
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// RateLimitMiddleware returns a gin middleware that rate limits by IP
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the key
		key := c.ClientIP()

		if !limiter.Allow(key) {
			response.Error(c, http.StatusTooManyRequests, "rate_limit_exceeded", "Too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUserMiddleware returns a gin middleware that rate limits by user ID
// Should be used after AuthMiddleware
func RateLimitByUserMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get user ID from context, fall back to IP
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			key = "user:" + userID.(string)
		}

		if !limiter.Allow(key) {
			response.Error(c, http.StatusTooManyRequests, "rate_limit_exceeded", "Too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}
