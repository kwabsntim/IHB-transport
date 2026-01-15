package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const EmailKey contextKey = "email"
const RoleKey contextKey = "role"

// authentication middleware //remember c.abort() stops the execution of further middleware/handlers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Safety check: ensure claims is not nil
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		userID, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract the role of the user
		role, ok := (*claims)["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store user info in Gin context for handlers to use
		c.Set("email", userID)
		c.Set("role", role)

		// Continue to next middleware/handler
		c.Next()
	}
}

// the middleware that extracts the role of the user from the response
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No role found in context"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot access this resource"})
			c.Abort()
			return
		}

		// Continue to next middleware/handler
		c.Next()
	}
}
