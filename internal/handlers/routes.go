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
		public.GET("/reviews/all", handler.GetAllReviewsHandler)
		public.GET("/reviews/:id", handler.GetReviewByIDHandler)
		public.DELETE("/reviews/:id/delete", handler.DeleteReviewHandler)

		
		// Track delivery by ID (clients)
		public.GET("/deliveries/:id", handler.GetDeliveryByIDHandler)

		// Get deliveries by email (clients checking their deliveries)
		public.GET("/deliveries/track", handler.GetDeliveriesByEmailHandler)

		// Client actions on instant quotes
		public.POST("/instant-quotes/:id/accept", handler.AcceptInstantQuotePriceHandler)
		public.GET("/instant-quotes/:id/accept", handler.AcceptInstantQuotePriceHandlerGET) // For email links
		public.POST("/instant-quotes/:id/decline", handler.DeclineInstantQuotePriceHandler)
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

			// Get all instant quotes
			admin.GET("/instant-quotes", handler.GetAllInstantQuotesHandler)

			// Get individual instant quote by ID
			admin.GET("/instant-quotes/:id", handler.GetInstantQuoteHandler)

			// Set price for instant quote
			admin.POST("/instant-quotes/:id/price", handler.SetInstantQuotePriceHandler)

			// Admin accepts delivery and sends confirmation email to client
			admin.POST("/deliveries/:id/accept", handler.AdminAcceptDeliveryHandler)
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
