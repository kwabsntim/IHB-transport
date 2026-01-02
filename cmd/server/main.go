package main

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/handlers"
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

	router := gin.Default()

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Setup all routes
	handlers.SetupRoutes(router)

	log.Println("🚀 Server running on :8080")
	router.Run(":8080")
}
