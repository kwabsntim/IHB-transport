package main

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/handlers"
	"ihb-transport/internal/repository"
	"ihb-transport/internal/services"
	"ihb-transport/migrations"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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

	// ==================== INITIALIZE SERVICES ====================
	emailService := services.NewEmailService(emailLogRepo)
	deliveryService := services.NewDeliveryService(
		deliveryRepo,
		statusLogRepo,
		emailLogRepo,
		emailService,
	)

	// ==================== INITIALIZE HANDLERS ====================
	handler := handlers.NewHandler(deliveryService)

	// ==================== SETUP ROUTER ====================
	router := gin.Default()

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Setup all routes with handler
	handlers.SetupRoutes(router, handler)

	log.Println("🚀 Server running on :8080")
	log.Println("📋 API Endpoints:")

	router.Run(":8080")
}
