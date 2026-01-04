package handlers

import (
	"ihb-transport/internal/auth"
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"ihb-transport/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler holds all service dependencies
type Handler struct {
	deliveryService services.DeliveryServiceInterface
}

// NewHandler creates a new handler with service dependencies
func NewHandler(deliveryService services.DeliveryServiceInterface) *Handler {
	return &Handler{
		deliveryService: deliveryService,
	}
}

// ==================== LOGIN HANDLER ====================

// input and response structs for login
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
		return
	}

	// Query database for admin
	var admin models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password with bcrypt
	if !auth.CheckPasswordHash(req.Password, admin.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := auth.CreateToken(admin.Email, "admin")
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

// ==================== DELIVERY HANDLERS ====================

// CreateDeliveryInput defines the input for creating a delivery request
type CreateDeliveryInput struct {
	ClientName      string  `json:"client_name" binding:"required"`
	ClientEmail     string  `json:"client_email" binding:"required,email"`
	PickupAddress   string  `json:"pickup_address" binding:"required"`
	DropoffAddress  string  `json:"dropoff_address" binding:"required"`
	ItemDescription string  `json:"item_description"`
	Weight          float64 `json:"weight"`
}

// CreateDeliveryHandler handles creation of new delivery requests (public)
func (h *Handler) CreateDeliveryHandler(c *gin.Context) {
	var input CreateDeliveryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create delivery struct
	delivery := models.DeliveryRequest{
		ClientName:      input.ClientName,
		ClientEmail:     input.ClientEmail,
		PickupAddress:   input.PickupAddress,
		DropoffAddress:  input.DropoffAddress,
		ItemDescription: input.ItemDescription,
		Weight:          input.Weight,
	}

	// Call service
	if err := h.deliveryService.CreateDeliveryRequest(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Delivery request created successfully",
		"delivery": delivery,
	})
}

// SetPriceInput defines the input for setting delivery price
type SetPriceInput struct {
	Price float64 `json:"price" binding:"required,gt=0"`
}

// SetDeliveryPriceHandler handles admin setting price for a delivery
func (h *Handler) SetDeliveryPriceHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	var input SetPriceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid price is required"})
		return
	}

	if err := h.deliveryService.SetDeliveryPrice(deliveryID, input.Price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Price set successfully"})
}

// AcceptDeliveryPriceHandler handles client accepting the price
func (h *Handler) AcceptDeliveryPriceHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.AcceptDeliveryPrice(deliveryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery price accepted successfully"})
}

// DeclineInput defines the input for declining delivery price
type DeclineInput struct {
	Reason string `json:"reason" binding:"required"`
}

// DeclineDeliveryPriceHandler handles client declining the price
func (h *Handler) DeclineDeliveryPriceHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	var input DeclineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decline reason is required"})
		return
	}

	if err := h.deliveryService.DeclineDeliveryPrice(deliveryID, input.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery price declined"})
}

// MarkAsPickedUpHandler handles driver marking delivery as picked up
func (h *Handler) MarkAsPickedUpHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.MarkAsPickedUp(deliveryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery marked as picked up"})
}

// MarkAsDeliveredHandler handles driver marking delivery as completed
func (h *Handler) MarkAsDeliveredHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.MarkAsDelivered(deliveryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery marked as completed"})
}

// GetAllDeliveriesHandler retrieves all deliveries (admin only)
func (h *Handler) GetAllDeliveriesHandler(c *gin.Context) {
	deliveries, err := h.deliveryService.GetAllDeliveries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":      len(deliveries),
		"deliveries": deliveries,
	})
}

// GetDeliveryByIDHandler retrieves a single delivery by ID
func (h *Handler) GetDeliveryByIDHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	delivery, err := h.deliveryService.GetDeliveryByID(deliveryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"delivery": delivery})
}

// GetDeliveriesByEmailHandler retrieves all deliveries for a client email
func (h *Handler) GetDeliveriesByEmailHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}

	deliveries, err := h.deliveryService.GetDeliveriesByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":      len(deliveries),
		"deliveries": deliveries,
	})
}

// GetDeliveriesByStatusHandler retrieves deliveries filtered by status
func (h *Handler) GetDeliveriesByStatusHandler(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status query parameter is required"})
		return
	}

	deliveries, err := h.deliveryService.GetDeliveriesByStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     status,
		"count":      len(deliveries),
		"deliveries": deliveries,
	})
}
