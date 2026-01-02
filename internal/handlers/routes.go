package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// Public routes
	router.POST("/login", LoginHandler)

}
