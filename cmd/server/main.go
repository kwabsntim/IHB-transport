package main

import (
	"IHB-transport/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

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
