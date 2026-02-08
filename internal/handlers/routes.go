package handlers

import (
	"ihb-transport/internal/auth"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// ==================== PUBLIC ROUTES ====================

	// Authentication
	router.POST("/login", LoginHandler)

	// Public delivery endpoints (anyone can create delivery request or track)
	public := router.Group("/api/public")
	{
		// Create delivery request (clients)
		public.POST("/deliveries", handler.CreateDeliveryHandler)

		// Create instant quote (clients)
		public.POST("/instant-quote", handler.CreateInstantQuoteHandler)

		// Reviews endpoints
		public.POST("/reviews", handler.CreateReviewHandler)
		public.GET("/reviews/:id", handler.GetReviewByIDHandler)

		
		// Track delivery by ID (clients)
		public.GET("/deliveries/:id", handler.GetDeliveryByIDHandler)

		// Get deliveries by email (clients checking their deliveries)
		public.GET("/deliveries/track", handler.GetDeliveriesByEmailHandler)

		// Client actions on their deliveries
		public.POST("/deliveries/:id/accept", handler.AcceptDeliveryPriceHandler)
		public.GET("/deliveries/:id/accept", handler.AcceptDeliveryPriceHandlerGET) // For email links
		
		public.POST("/deliveries/:id/decline", handler.DeclineDeliveryPriceHandler)
	}

	// ==================== PROTECTED ROUTES (Admin/Driver) ====================

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware()) // Requires valid JWT token
	{
		// Admin routes - manage all deliveries
		admin := protected.Group("/admin")
		admin.Use(auth.RoleMiddleware("admin")) // Requires admin role
		{
			// Get all deliveries
			admin.GET("/deliveries", handler.GetAllDeliveriesHandler)

			// Get deliveries by status
			admin.GET("/deliveries/status", handler.GetDeliveriesByStatusHandler)

			// Set price for delivery
			admin.POST("/deliveries/:id/price", handler.SetDeliveryPriceHandler)

			// Get all instant quotes
			admin.GET("/instant-quotes", handler.GetAllInstantQuotesHandler)

			// Get individual instant quote by ID
			admin.GET("/instant-quotes/:id", handler.GetInstantQuoteHandler)
		}

		// Driver routes - update delivery status
		driver := protected.Group("/driver")
		{
			// Mark as picked up
			driver.POST("/deliveries/:id/pickup", handler.MarkAsPickedUpHandler)

			// Mark as delivered
			driver.POST("/deliveries/:id/complete", handler.MarkAsDeliveredHandler)

			// Get deliveries assigned to driver (optional - could filter by status)
			driver.GET("/deliveries", handler.GetDeliveriesByStatusHandler)
		}
	}
}
