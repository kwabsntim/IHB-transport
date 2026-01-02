package auth

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// input and reponse structs for login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message,omitempty"`
}

// LoginHandler handles user login requests

func LoginHandler(c *gin.Context) {
	var req LoginInput
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
	}

	// Query database for admin
	var admin models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password with bcrypt
	if !CheckPasswordHash(req.Password, admin.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := CreateToken(admin.Email, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	response := LoginResponse{
		Token:   token,
		Message: "Login successful",
	}
	c.JSON(http.StatusOK, response)

}
