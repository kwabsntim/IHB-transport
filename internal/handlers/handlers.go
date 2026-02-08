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
	reviewService   services.ReviewServiceInterface
}

// NewHandler creates a new handler with service dependencies
func NewHandler(deliveryService services.DeliveryServiceInterface, reviewService services.ReviewServiceInterface) *Handler {
	return &Handler{
		deliveryService: deliveryService,
		reviewService:   reviewService,
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
	PriceOption     string `json:"price_option"`
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
		PriceOption:     input.PriceOption,
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

// AdminAcceptDeliveryHandler handles admin accepting the delivery (sends confirmation email to client)
func (h *Handler) AdminAcceptDeliveryHandler(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.AdminAcceptDeliveryPrice(deliveryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery accepted by admin and confirmation email sent to client"})
}

// AcceptDeliveryPriceHandlerGET handles client accepting the price from email link (GET - HTML response)
func (h *Handler) AcceptDeliveryPriceHandlerGET(c *gin.Context) {
	deliveryID := c.Param("id")

	if err := h.deliveryService.AcceptDeliveryPrice(deliveryID); err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Cannot Accept Quote - IHB Transport</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<style>
					body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background-color: #f5f5f5; padding: 20px; }
					.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; width: 100%; }
					.error-icon { font-size: 60px; margin-bottom: 20px; }
					h1 { color: #dc3545; margin-bottom: 20px; font-size: 24px; }
					.error-message { color: #333; line-height: 1.6; background: #f8d7da; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #dc3545; }
					p { color: #666; line-height: 1.6; }
					.delivery-id { background: #f8f9fa; padding: 10px; border-radius: 5px; margin: 15px 0; font-family: monospace; font-size: 12px; word-break: break-all; }
					.contact-box { background: #e7f3ff; padding: 15px; border-radius: 5px; margin-top: 20px; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="error-icon">⚠️</div>
					<h1>Cannot Accept Quote</h1>
					<div class="error-message">`+err.Error()+`</div>
					<div class="delivery-id">Delivery ID: `+deliveryID+`</div>
					<div class="contact-box">
						<p style="margin: 0; font-weight: bold; color: #0056b3;">Need Help?</p>
						<p style="margin: 5px 0 0 0; font-size: 14px;">Contact IHB Transport support for assistance with this delivery.</p>
					</div>
					<p style="margin-top: 20px; font-size: 12px; color: #999;">You can close this window now.</p>
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
				<title>Cannot Decline Quote - IHB Transport</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<style>
					body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background-color: #f5f5f5; padding: 20px; }
					.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; width: 100%; }
					.error-icon { font-size: 60px; margin-bottom: 20px; }
					h1 { color: #dc3545; margin-bottom: 20px; font-size: 24px; }
					.error-message { color: #333; line-height: 1.6; background: #f8d7da; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #dc3545; }
					p { color: #666; line-height: 1.6; }
					.delivery-id { background: #f8f9fa; padding: 10px; border-radius: 5px; margin: 15px 0; font-family: monospace; font-size: 12px; word-break: break-all; }
					.contact-box { background: #e7f3ff; padding: 15px; border-radius: 5px; margin-top: 20px; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="error-icon">⚠️</div>
					<h1>Cannot Decline Quote</h1>
					<div class="error-message">`+err.Error()+`</div>
					<div class="delivery-id">Delivery ID: `+deliveryID+`</div>
					<div class="contact-box">
						<p style="margin: 0; font-weight: bold; color: #0056b3;">Need Help?</p>
						<p style="margin: 5px 0 0 0; font-size: 14px;">Contact IHB Transport support for assistance with this delivery.</p>
					</div>
					<p style="margin-top: 20px; font-size: 12px; color: #999;">You can close this window now.</p>
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

// CreateInstantQuoteHandler creates and persists an instant quote and optionally triggers an email
func (h *Handler) CreateInstantQuoteHandler(c *gin.Context) {
	type InstantQuoteInput struct {
		PickupPoint     string `json:"pickup_point" binding:"required"`
		DeliveryAddress string `json:"delivery_address" binding:"required"`
		Weight          string `json:"weight" binding:"required"`
		ClientEmail     string `json:"client_email" binding:"required,email"`
	}

	var input InstantQuoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote := models.InstantQuote{
		PickupPoint:     input.PickupPoint,
		DeliveryAddress: input.DeliveryAddress,
		Weight:          input.Weight,
		ClientEmail:     input.ClientEmail,
	}

	if err := h.deliveryService.GetInstantQuote(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Instant quote created", "instant_quote": quote})
}

// CreateReviewHandler creates a new review (public)
func (h *Handler) CreateReviewHandler(c *gin.Context) {
	var input struct {
		ClientName string `json:"client_name" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := models.Reviews{
		ClientName: input.ClientName,
		Content:    input.Content,
	}

	if err := h.reviewService.CreateReview(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review created", "review": review})
}

// GetAllReviewsHandler retrieves all reviews
func (h *Handler) GetAllReviewsHandler(c *gin.Context) {
	reviews, err := h.reviewService.GetAllReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(reviews),
		"reviews": reviews,
	})
}

// GetReviewByIDHandler retrieves a review by ID
func (h *Handler) GetReviewByIDHandler(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "review id is required"})
		return
	}

	rev, err := h.reviewService.GetReviewByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": rev})
}

// GetAllInstantQuotesHandler retrieves all instant quotes (admin only)
func (h *Handler) GetAllInstantQuotesHandler(c *gin.Context) {
	var quotes []models.InstantQuote
	if err := database.DB.Order("created_at DESC").Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch instant quotes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":          len(quotes),
		"instant_quotes": quotes,
	})
}

// GetInstantQuoteHandler retrieves an instant quote by ID
func (h *Handler) GetInstantQuoteHandler(c *gin.Context) {
	quoteID := c.Param("id")
	if quoteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instant quote id is required"})
		return
	}

	var quote models.InstantQuote
	if err := database.DB.First(&quote, "id = ?", quoteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "instant quote not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"instant_quote": quote})
}

// SetInstantQuotePriceHandler handles admin setting price for an instant quote
func (h *Handler) SetInstantQuotePriceHandler(c *gin.Context) {
	quoteID := c.Param("id")

	var input SetPriceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid price is required"})
		return
	}

	if err := h.deliveryService.SetInstantQuotePrice(quoteID, input.Price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instant quote price set successfully and email sent to client"})
}

// AcceptInstantQuotePriceHandler handles client accepting the instant quote price (POST - JSON response)
func (h *Handler) AcceptInstantQuotePriceHandler(c *gin.Context) {
	quoteID := c.Param("id")

	if err := h.deliveryService.AcceptInstantQuotePrice(quoteID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instant quote price accepted successfully"})
}

// AcceptInstantQuotePriceHandlerGET handles client accepting the instant quote price from email link (GET - HTML response)
func (h *Handler) AcceptInstantQuotePriceHandlerGET(c *gin.Context) {
	quoteID := c.Param("id")

	if err := h.deliveryService.AcceptInstantQuotePrice(quoteID); err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Cannot Accept Quote - IHB Transport</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<style>
					body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background-color: #f5f5f5; padding: 20px; }
					.container { text-align: center; padding: 40px; background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 500px; width: 100%; }
					.error-icon { font-size: 60px; margin-bottom: 20px; }
					h1 { color: #dc3545; margin-bottom: 20px; font-size: 24px; }
					.error-message { color: #333; line-height: 1.6; background: #f8d7da; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #dc3545; }
					p { color: #666; line-height: 1.6; }
					.quote-id { background: #f8f9fa; padding: 10px; border-radius: 5px; margin: 15px 0; font-family: monospace; font-size: 12px; word-break: break-all; }
					.contact-box { background: #e7f3ff; padding: 15px; border-radius: 5px; margin-top: 20px; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="error-icon">⚠️</div>
					<h1>Cannot Accept Quote</h1>
					<div class="error-message">`+err.Error()+`</div>
					<div class="quote-id">Quote ID: `+quoteID+`</div>
					<div class="contact-box">
						<p style="margin: 0; font-weight: bold; color: #0056b3;">Need Help?</p>
						<p style="margin: 5px 0 0 0; font-size: 14px;">Contact IHB Transport support for assistance.</p>
					</div>
					<p style="margin-top: 20px; font-size: 12px; color: #999;">You can close this window now.</p>
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
				.quote-id { background: #e7f3ff; padding: 10px; border-radius: 5px; margin: 20px 0; font-family: monospace; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="success-icon">✅</div>
				<h1>Quote Accepted!</h1>
				<p>Thank you for accepting our quote. We will be in touch shortly.</p>
				<div class="quote-id">Quote ID: `+quoteID+`</div>
				<p>You can close this window now.</p>
			</div>
		</body>
		</html>
	`))
}

// DeclineInstantQuotePriceHandler handles client declining the instant quote price
func (h *Handler) DeclineInstantQuotePriceHandler(c *gin.Context) {
	quoteID := c.Param("id")

	var input struct {
		Reason string `json:"reason"`
	}
	// Reason is optional, so we don't require binding
	c.ShouldBindJSON(&input)

	if err := h.deliveryService.DeclineInstantQuotePrice(quoteID, input.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instant quote declined"})
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
