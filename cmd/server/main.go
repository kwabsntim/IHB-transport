package main

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/handlers"
	"ihb-transport/internal/repository"
	"ihb-transport/internal/services"
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
	log.Println("   Public:")
	log.Println("     POST   /login")
	log.Println("     POST   /api/public/deliveries")
	log.Println("     GET    /api/public/deliveries/:id")
	log.Println("     GET    /api/public/deliveries/track?email=...")
	log.Println("     POST   /api/public/deliveries/:id/accept")
	log.Println("     POST   /api/public/deliveries/:id/decline")
	log.Println("   Admin (requires auth):")
	log.Println("     GET    /api/admin/deliveries")
	log.Println("     GET    /api/admin/deliveries/status?status=...")
	log.Println("     POST   /api/admin/deliveries/:id/price")
	log.Println("   Driver (requires auth):")
	log.Println("     POST   /api/driver/deliveries/:id/pickup")
	log.Println("     POST   /api/driver/deliveries/:id/complete")

	router.Run(":8080")
}
