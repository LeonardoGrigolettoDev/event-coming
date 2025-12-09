package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsCollector collects HTTP request metrics
type MetricsCollector struct {
	mu sync.RWMutex

	// Request counts by method and status
	RequestsTotal map[string]map[int]int64

	// Request latencies by endpoint (in milliseconds)
	RequestLatencies map[string][]float64

	// Error counts
	ErrorsTotal int64

	// Active requests
	ActiveRequests int64
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		RequestsTotal:    make(map[string]map[int]int64),
		RequestLatencies: make(map[string][]float64),
	}
}

// MetricsMiddleware collects request metrics
func MetricsMiddleware(collector *MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment active requests
		collector.mu.Lock()
		collector.ActiveRequests++
		collector.mu.Unlock()

		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := float64(time.Since(start).Milliseconds())
		status := c.Writer.Status()
		method := c.Request.Method
		key := method + " " + path

		// Update metrics
		collector.mu.Lock()
		defer collector.mu.Unlock()

		// Decrement active requests
		collector.ActiveRequests--

		// Increment request count
		if collector.RequestsTotal[key] == nil {
			collector.RequestsTotal[key] = make(map[int]int64)
		}
		collector.RequestsTotal[key][status]++

		// Record latency (keep last 1000 samples per endpoint)
		collector.RequestLatencies[key] = append(collector.RequestLatencies[key], latency)
		if len(collector.RequestLatencies[key]) > 1000 {
			collector.RequestLatencies[key] = collector.RequestLatencies[key][1:]
		}

		// Count errors (5xx)
		if status >= 500 {
			collector.ErrorsTotal++
		}
	}
}

// GetMetrics returns current metrics as a map (for JSON response)
func (c *MetricsCollector) GetMetrics() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Calculate summary stats for latencies
	latencySummaries := make(map[string]map[string]float64)
	for endpoint, latencies := range c.RequestLatencies {
		if len(latencies) == 0 {
			continue
		}

		summary := map[string]float64{
			"count": float64(len(latencies)),
			"avg":   average(latencies),
			"min":   min(latencies),
			"max":   max(latencies),
			"p50":   percentile(latencies, 50),
			"p95":   percentile(latencies, 95),
			"p99":   percentile(latencies, 99),
		}
		latencySummaries[endpoint] = summary
	}

	// Build request counts
	requestCounts := make(map[string]map[string]int64)
	for endpoint, statusCounts := range c.RequestsTotal {
		requestCounts[endpoint] = make(map[string]int64)
		for status, count := range statusCounts {
			requestCounts[endpoint][strconv.Itoa(status)] = count
		}
	}

	return map[string]interface{}{
		"active_requests": c.ActiveRequests,
		"errors_total":    c.ErrorsTotal,
		"requests_total":  requestCounts,
		"latency_ms":      latencySummaries,
	}
}

// PrometheusFormat returns metrics in Prometheus text format
func (c *MetricsCollector) PrometheusFormat() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var output string

	// Active requests
	output += "# HELP http_active_requests Number of currently active HTTP requests\n"
	output += "# TYPE http_active_requests gauge\n"
	output += "http_active_requests " + strconv.FormatInt(c.ActiveRequests, 10) + "\n\n"

	// Total errors
	output += "# HELP http_errors_total Total number of HTTP 5xx errors\n"
	output += "# TYPE http_errors_total counter\n"
	output += "http_errors_total " + strconv.FormatInt(c.ErrorsTotal, 10) + "\n\n"

	// Request counts by status
	output += "# HELP http_requests_total Total number of HTTP requests\n"
	output += "# TYPE http_requests_total counter\n"
	for endpoint, statusCounts := range c.RequestsTotal {
		for status, count := range statusCounts {
			output += "http_requests_total{endpoint=\"" + endpoint + "\",status=\"" + strconv.Itoa(status) + "\"} " + strconv.FormatInt(count, 10) + "\n"
		}
	}

	output += "\n"

	// Latency histograms (simplified - just summary stats)
	output += "# HELP http_request_duration_ms HTTP request latency in milliseconds\n"
	output += "# TYPE http_request_duration_ms summary\n"
	for endpoint, latencies := range c.RequestLatencies {
		if len(latencies) == 0 {
			continue
		}
		output += "http_request_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"0.5\"} " + strconv.FormatFloat(percentile(latencies, 50), 'f', 2, 64) + "\n"
		output += "http_request_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"0.95\"} " + strconv.FormatFloat(percentile(latencies, 95), 'f', 2, 64) + "\n"
		output += "http_request_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"0.99\"} " + strconv.FormatFloat(percentile(latencies, 99), 'f', 2, 64) + "\n"
		output += "http_request_duration_ms_sum{endpoint=\"" + endpoint + "\"} " + strconv.FormatFloat(sum(latencies), 'f', 2, 64) + "\n"
		output += "http_request_duration_ms_count{endpoint=\"" + endpoint + "\"} " + strconv.Itoa(len(latencies)) + "\n"
	}

	return output
}

// Helper functions for statistics
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return sum(values) / float64(len(values))
}

func sum(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Simple percentile calculation (not sorted, approximation)
	// For production, use a proper algorithm or library
	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort for small datasets
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	index := int(float64(len(sorted)-1) * p / 100)
	return sorted[index]
}
