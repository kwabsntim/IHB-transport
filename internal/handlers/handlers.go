package handlers

import (
	"ihb-transport/internal/auth"
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"ihb-transport/internal/services"
	"net/http"
	"time"

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
	ClientName  string `json:"client_name" binding:"required"`
	ClientEmail string `json:"client_email" binding:"required,email"`

	// Pickup Address Details
	PickupStreet   string `json:"pickup_street" binding:"required"`
	PickupCity     string `json:"pickup_city" binding:"required"`
	PickupPostCode string `json:"pickup_post_code" binding:"required"`
	PickupCountry  string `json:"pickup_country" binding:"required"`

	// Dropoff Address Details
	DropoffStreet   string `json:"dropoff_street" binding:"required"`
	DropoffCity     string `json:"dropoff_city" binding:"required"`
	DropoffPostCode string `json:"dropoff_post_code" binding:"required"`
	DropoffCountry  string `json:"dropoff_country" binding:"required"`

	ItemDescription string `json:"item_description"`
	Items           string `json:"items" binding:"required"`
	Weight          string `json:"weight"`
	Service         string `json:"service" binding:"required"`
	PickupDate      string `json:"pickup_date" binding:"required"` // Format: YYYY-MM-DD
}

// CreateDeliveryHandler handles creation of new delivery requests (public)
func (h *Handler) CreateDeliveryHandler(c *gin.Context) {
	var input CreateDeliveryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse pickup date
	pickupDate, err := time.Parse("2006-01-02", input.PickupDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pickup_date format. Use YYYY-MM-DD"})
		return
	}

	// Create delivery struct
	delivery := models.DeliveryRequest{
		ClientName:      input.ClientName,
		ClientEmail:     input.ClientEmail,
		PickupStreet:    input.PickupStreet,
		PickupCity:      input.PickupCity,
		PickupPostCode:  input.PickupPostCode,
		PickupCountry:   input.PickupCountry,
		DropoffStreet:   input.DropoffStreet,
		DropoffCity:     input.DropoffCity,
		DropoffPostCode: input.DropoffPostCode,
		DropoffCountry:  input.DropoffCountry,
		ItemDescription: input.ItemDescription,
		Items:           input.Items,
		Weight:          input.Weight,
		Service:         input.Service,
		PickupDate:      &pickupDate,
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

// AcceptDeliveryPriceHandler handles client accepting the price (POST - JSON response)
func (h *Handler) AcceptDeliveryPriceHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.AcceptDeliveryPrice(deliveryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery price accepted successfully"})
}

// AcceptDeliveryPriceHandlerGET handles client accepting the price from email link (GET - HTML response)
func (h *Handler) AcceptDeliveryPriceHandlerGET(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.AcceptDeliveryPrice(deliveryID); err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Error - IHB Transport</title>
				<style>
					body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
					.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; }
					.error-icon { font-size: 60px; margin-bottom: 20px; }
					h1 { color: #dc3545; margin-bottom: 20px; }
					p { color: #666; line-height: 1.6; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="error-icon">❌</div>
					<h1>Unable to Accept Quote</h1>
					<p>`+err.Error()+`</p>
					<p style="margin-top: 30px; font-size: 14px; color: #999;">If you need assistance, please contact IHB Transport support.</p>
				</div>
			</body>
			</html>
		`))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Quote Accepted - IHB Transport</title>
			<style>
				body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
				.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; }
				.success-icon { font-size: 60px; margin-bottom: 20px; }
				h1 { color: #28a745; margin-bottom: 20px; }
				p { color: #666; line-height: 1.6; }
				.delivery-id { background: #e7f3ff; padding: 10px; border-radius: 5px; margin: 20px 0; font-family: monospace; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="success-icon">✅</div>
				<h1>Quote Accepted Successfully!</h1>
				<p>Thank you for accepting our quote. We've received your confirmation.</p>
				<div class="delivery-id">Delivery ID: `+deliveryID+`</div>
				<p><strong>What happens next?</strong></p>
				<p>Our team will review your acceptance and assign a driver shortly. You'll receive an email confirmation once a driver is on the way.</p>
				<p style="margin-top: 30px; font-size: 14px; color: #999;">You can close this window now.</p>
			</div>
		</body>
		</html>
	`))
}

// DeclineInput defines the input for declining delivery price
type DeclineInput struct {
	Reason string `json:"reason" binding:"required"`
}

// DeclineDeliveryPriceHandler handles client declining the price (POST - JSON response)
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

// DeclineDeliveryPriceHandlerGET handles client declining the price from email link (GET - HTML response)
func (h *Handler) DeclineDeliveryPriceHandlerGET(c *gin.Context) {
	deliveryID := c.Param("id")
	reason := c.Query("reason") // Optional reason from query parameter

	if reason == "" {
		reason = "No reason provided"
	}

	if err := h.deliveryService.DeclineDeliveryPrice(deliveryID, reason); err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Error - IHB Transport</title>
				<style>
					body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
					.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; }
					.error-icon { font-size: 60px; margin-bottom: 20px; }
					h1 { color: #dc3545; margin-bottom: 20px; }
					p { color: #666; line-height: 1.6; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="error-icon">❌</div>
					<h1>Unable to Decline Quote</h1>
					<p>`+err.Error()+`</p>
					<p style="margin-top: 30px; font-size: 14px; color: #999;">If you need assistance, please contact IHB Transport support.</p>
				</div>
			</body>
			</html>
		`))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Quote Declined - IHB Transport</title>
			<style>
				body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
				.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; }
				.info-icon { font-size: 60px; margin-bottom: 20px; }
				h1 { color: #666; margin-bottom: 20px; }
				p { color: #666; line-height: 1.6; }
				.delivery-id { background: #f8f9fa; padding: 10px; border-radius: 5px; margin: 20px 0; font-family: monospace; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="info-icon">📋</div>
				<h1>Quote Declined</h1>
				<p>We've received your decision to decline the quote.</p>
				<div class="delivery-id">Delivery ID: `+deliveryID+`</div>
				<p>Thank you for considering IHB Transport. If you'd like to discuss alternative options or have questions about the quote, please feel free to contact us.</p>
				<p style="margin-top: 30px; font-size: 14px; color: #999;">You can close this window now.</p>
			</div>
		</body>
		</html>
	`))
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
