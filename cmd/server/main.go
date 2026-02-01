package main

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/handlers"
	"ihb-transport/internal/repository"
	"ihb-transport/internal/services"
	"ihb-transport/migrations"
	"ihb-transport/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional for production)
	_ = godotenv.Load() // Ignore error - .env is optional in production

	// Connect to Postgres
	database.Connection()

	// Run migrations automatically (safe to run multiple times)
	log.Println("🔄 Running database migrations...")
	migrations.Migrate()
	log.Println("✅ Migrations complete")

	// ==================== INITIALIZE REPOSITORIES ====================
	deliveryRepo := repository.NewDeliveryRepository()
	statusLogRepo := repository.NewStatusLogRepository()
	emailLogRepo := repository.NewEmailLogRepository()
	reviewsRepo := repository.NewReviewsRepository()

	// ==================== INITIALIZE SERVICES ====================
	emailService := services.NewEmailService(emailLogRepo)
	deliveryService := services.NewDeliveryService(
		deliveryRepo,
		statusLogRepo,
		emailLogRepo,
		emailService,
	)
	reviewsService := services.NewReviewService(reviewsRepo)

	// ==================== INITIALIZE HANDLERS ====================
	handler := handlers.NewHandler(deliveryService, reviewsService)

	// ==================== SETUP ROUTER ====================
	router := gin.New() // Use gin.New() instead of gin.Default() for custom middleware

	// ==================== APPLY GLOBAL MIDDLEWARE ====================

	// 1. Recovery middleware (handles panics)
	router.Use(utils.RecoveryMiddleware())

	// 2. Request logger (logs all requests)
	router.Use(utils.RequestLoggerMiddleware())

	// 3. CORS middleware (allows cross-origin requests)
	router.Use(utils.CORSMiddleware())

	// 4. Security headers
	router.Use(utils.SecurityHeadersMiddleware())

	// 5. Rate limiting (10 requests per second per IP, burst of 20)
	utils.InitRateLimiter(10, 20)
	router.Use(utils.RateLimitMiddleware())

	// 6. Request size limit (10MB)
	router.Use(utils.RequestSizeLimitMiddleware(10 * 1024 * 1024))

	// 7. Request timeout (30 seconds)
	// Note: Commented out as it may interfere with long-running operations
	// router.Use(utils.TimeoutMiddleware(30 * time.Second))

	// 8. Error handler (catches and formats errors)
	router.Use(utils.ErrorHandlerMiddleware())

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "healthy",
		})
	})

	// Setup all routes with handler
	handlers.SetupRoutes(router, handler)

	log.Println("🚀 Server running on :8080")
	log.Println("📋 API Endpoints:")

	router.Run(":8080")
}
