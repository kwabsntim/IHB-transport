package main

import (
	"ihb-transport/internal/database"
	"net/http"

	"log"

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
	router.GET("/ping", func(c *gin.Context) {
		// Respond with a JSON message and an HTTP status 200 OK
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080")

}
