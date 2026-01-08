package utils

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ==================== CORS MIDDLEWARE ====================

// CORSMiddleware handles Cross-Origin Resource Sharing
// Allows frontend applications from different domains to access the API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin (configure this based on your needs)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// Allow specific headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")

		// Allow credentials (cookies, authorization headers)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Cache preflight requests for 12 hours
		c.Writer.Header().Set("Access-Control-Max-Age", "43200")

		// Expose specific headers to the client
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSMiddlewareStrict provides more restrictive CORS configuration
// Use this in production with specific allowed origins
func CORSMiddlewareStrict(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Max-Age", "43200")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ==================== REQUEST LOGGER MIDDLEWARE ====================

// RequestLoggerMiddleware logs all incoming HTTP requests with detailed information
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request details
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Get error if any
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log format: [timestamp] method path | status | duration | IP | user-agent
		logMessage := fmt.Sprintf(
			"[%s] %s %s",
			startTime.Format("2006-01-02 15:04:05"),
			method,
			path,
		)

		if query != "" {
			logMessage += fmt.Sprintf("?%s", query)
		}

		logMessage += fmt.Sprintf(
			" | Status: %d | Duration: %v | IP: %s | UA: %s",
			statusCode,
			duration,
			clientIP,
			userAgent,
		)

		if errorMessage != "" {
			logMessage += fmt.Sprintf(" | Errors: %s", errorMessage)
		}

		// Color-code based on status
		if statusCode >= 500 {
			fmt.Printf("🔴 %s\n", logMessage)
		} else if statusCode >= 400 {
			fmt.Printf("🟡 %s\n", logMessage)
		} else if statusCode >= 300 {
			fmt.Printf("🟢 %s\n", logMessage)
		} else {
			fmt.Printf("✅ %s\n", logMessage)
		}
	}
}

// ==================== RATE LIMITER MIDDLEWARE ====================

// RateLimiter stores rate limiters for each IP address
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// rate: requests per second
// burst: maximum burst size
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter returns the rate limiter for a given IP
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// CleanupOldLimiters removes inactive limiters (call periodically)
func (rl *RateLimiter) CleanupOldLimiters() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Clear all limiters (in production, implement smarter cleanup based on last access time)
	rl.limiters = make(map[string]*rate.Limiter)
}

// Global rate limiter instance
var globalRateLimiter *RateLimiter

// InitRateLimiter initializes the global rate limiter
func InitRateLimiter(requestsPerSecond float64, burstSize int) {
	globalRateLimiter = NewRateLimiter(rate.Limit(requestsPerSecond), burstSize)

	// Start cleanup goroutine (runs every 5 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			globalRateLimiter.CleanupOldLimiters()
		}
	}()
}

// RateLimitMiddleware limits the number of requests per IP address
// Default: 10 requests per second with burst of 20
func RateLimitMiddleware() gin.HandlerFunc {
	// Initialize with default values if not already initialized
	if globalRateLimiter == nil {
		InitRateLimiter(10, 20) // 10 req/sec, burst of 20
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := globalRateLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": "1 second",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewareCustom creates a custom rate limit middleware
func RateLimitMiddlewareCustom(requestsPerSecond float64, burstSize int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(requestsPerSecond), burstSize)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		ipLimiter := limiter.GetLimiter(ip)

		if !ipLimiter.Allow() {
			retryAfter := time.Second * time.Duration(1.0/requestsPerSecond)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Too many requests. Limit: %.0f requests per second.", requestsPerSecond),
				"retry_after": retryAfter.String(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== ERROR HANDLER MIDDLEWARE ====================

// ErrorHandlerMiddleware catches and handles all errors in a consistent format
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute the request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			// Determine status code
			statusCode := c.Writer.Status()
			if statusCode == http.StatusOK {
				// Default to 500 if status wasn't set
				statusCode = http.StatusInternalServerError
			}

			// Create error response
			errorResponse := gin.H{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
			}

			// Add more details for development (comment out in production)
			if gin.Mode() == gin.DebugMode {
				errorResponse["method"] = c.Request.Method
				errorResponse["timestamp"] = time.Now().Format(time.RFC3339)
			}

			// Log the error
			fmt.Printf("❌ Error on %s %s: %s\n", c.Request.Method, c.Request.URL.Path, err.Error())

			// Send JSON response
			c.JSON(statusCode, errorResponse)
			return
		}
	}
}

// RecoveryMiddleware handles panics and converts them to 500 errors
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				fmt.Printf("🚨 PANIC RECOVERED: %v\n", err)
				fmt.Printf("   Path: %s %s\n", c.Request.Method, c.Request.URL.Path)
				fmt.Printf("   IP: %s\n", c.ClientIP())

				// Send error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred. Please try again later.",
				})

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ==================== SECURITY MIDDLEWARE ====================

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Enforce HTTPS (uncomment in production with HTTPS)
		// c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Control referrer information
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (adjust based on your needs)
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// ==================== REQUEST SIZE LIMIT MIDDLEWARE ====================

// RequestSizeLimitMiddleware limits the size of request bodies
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set max request body size (default: 10MB)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}

// ==================== TIMEOUT MIDDLEWARE ====================

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a channel to signal completion
		done := make(chan struct{})

		// Run the request in a goroutine
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Request completed
			return
		case <-time.After(timeout):
			// Timeout reached
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Request timeout",
				"message": fmt.Sprintf("Request took longer than %v", timeout),
			})
			c.Abort()
			return
		}
	}
}
