package services

import (
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
	"ihb-transport/internal/utils"
	"net/smtp"
	"os"

	"github.com/google/uuid"
)

type emailService struct {
	emailLogRepo repository.EmailLogInterface
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
	baseURL      string
	enabled      bool
}

// NewEmailService creates a new email service
func NewEmailService(emailLogRepo repository.EmailLogInterface) EmailServiceInterface {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")

	// Check if SMTP is enabled
	enabled := smtpHost != "" && smtpPort != "" && smtpUser != "" && smtpPassword != ""

	if !enabled {
		fmt.Println("⚠️  SMTP not configured - emails will be logged to console only")
		fmt.Println("   To enable SMTP, set: SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD")
	} else {
		fmt.Printf("✅ SMTP enabled - using %s:%s\n", smtpHost, smtpPort)
	}

	if fromEmail == "" {
		fromEmail = "noreply@ihbtransport.com"
	}
	if fromName == "" {
		fromName = "IHB Transport"
	}
	
	// Get base URL for email links (defaults to localhost for development)
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &emailService{
		emailLogRepo: emailLogRepo,
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
		fromName:     fromName,
		baseURL:      baseURL,
		enabled:      enabled,
	}
}

// SendRequestReceivedEmail sends confirmation email when request is created
func (s *emailService) SendRequestReceivedEmail(clientEmail, clientName, deliveryID, pickupAddress, dropoffAddress, service, pickupDate string) error {
	subject := fmt.Sprintf("Delivery Request Received #%s", deliveryID)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 20px; }
				.info-box { background-color: white; padding: 15px; margin: 10px 0; border-left: 4px solid #4CAF50; }
				.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>✅ Request Received</h1>
				</div>
				<div class="content">
					<p>Dear %s,</p>
					<p>Thank you for choosing IHB Transport! We have received your delivery request.</p>
					
					<div class="info-box">
						<strong>Request ID:</strong> %s<br>
						<strong>Service:</strong> %s<br>
						<strong>Pickup Date:</strong> %s<br>
						<strong>Pickup Location:</strong> %s<br>
						<strong>Dropoff Location:</strong> %s
					</div>
					
					<p>Our team will review your request and send you a price quote shortly.</p>
					<p>You will receive another email once we have prepared your quote.</p>
				</div>
				<div class="footer">
					<p>IHB Transport - Reliable Delivery Services</p>
				</div>
			</div>
		</body>
		</html>
	`, clientName, deliveryID, service, pickupDate, pickupAddress, dropoffAddress)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "REQUEST_RECEIVED", err)
}

// SendPriceEmail sends email with quoted price
func (s *emailService) SendPriceEmail(clientEmail string, price float64, deliveryID string) error {
	subject := fmt.Sprintf("Price Quote for Delivery #%s", deliveryID)
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #007bff; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 20px; }
				.price-box { background-color: white; padding: 20px; margin: 20px 0; border-left: 4px solid #007bff; text-align: center; }
				.button-container { text-align: center; margin: 30px 0; }
				.button { display: inline-block; padding: 12px 30px; margin: 10px; text-decoration: none; border-radius: 5px; font-weight: bold; }
				.accept-btn { background-color: #28a745; color: white; }
				.decline-btn { background-color: #dc3545; color: white; }
				.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>💰 Your Quote is Ready!</h1>
				</div>
				<div class="content">
					<p>We're pleased to provide you with a quote for your delivery request.</p>
					
					<div class="price-box">
						<p><strong>Delivery ID:</strong> %s</p>
						<h2 style="margin: 10px 0; color: #007bff;">Price: %.2f DKK</h2>
					</div>
					
					<p>Please review the quote and choose an option:</p>
					
					<div class="button-container">
						<a href="%s/api/public/deliveries/%s/accept" class="button accept-btn">✅ Accept Quote</a>
						<a href="%s/api/public/deliveries/%s/decline" class="button decline-btn">❌ Decline Quote</a>
					</div>
					
					<p style="font-size: 12px; color: #666;">This quote is valid for 48 hours.</p>
				</div>
				<div class="footer">
					<p>IHB Transport - Reliable Delivery Services</p>
				</div>
			</div>
		</body>
		</html>
	`, deliveryID, price, s.baseURL, deliveryID, s.baseURL, deliveryID)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "PRICE_SENT", err)
}

// SendAcceptedEmail notifies that client accepted the price
func (s *emailService) SendAcceptedEmail(clientEmail, deliveryID string) error {
	subject := fmt.Sprintf("Delivery Confirmed #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2>🎉 Your Delivery is Confirmed!</h2>
		<p>Thank you for accepting our quote! Your delivery has been confirmed and scheduled.</p>
		<p><strong>Delivery ID:</strong> %s</p>
		<p><strong>Next Steps:</strong></p>
		<ol>
			<li>Our driver will pick up your package soon</li>
			<li>You'll receive an update when pickup is complete</li>
			<li>Track your delivery status in real-time</li>
		</ol>
		<p>We'll send you notifications as your delivery progresses.</p>
		<br>
		<p>Thank you for choosing IHB Transport!<br>IHB Transport Team</p>
	`, deliveryID)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "ACCEPTED", err)
}

// SendDeclinedEmail sends notification when client declines the price
func (s *emailService) SendDeclinedEmail(clientEmail, deliveryID, reason string) error {
	// Send notification to admin
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		adminSubject := fmt.Sprintf("Price Declined - Delivery #%s", deliveryID)
		adminBody := fmt.Sprintf(`
			<h2>Price Declined by Client</h2>
			<p><strong>Delivery ID:</strong> %s</p>
			<p><strong>Client Email:</strong> %s</p>
			<p><strong>Decline Reason:</strong> %s</p>
			<p>You may want to follow up with the client with a revised quote.</p>
		`, deliveryID, clientEmail, reason)
		s.sendEmail(adminEmail, adminSubject, adminBody)
	}

	// Send confirmation to client
	subject := fmt.Sprintf("Delivery Request Declined #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2>Request Declined</h2>
		<p>We've received your decision to decline the delivery quote.</p>
		<p><strong>Delivery ID:</strong> %s</p>
		<p>We understand that pricing is important, and we appreciate your honesty.</p>
		<p>If you'd like to discuss alternative options or have any questions, please don't hesitate to contact us.</p>
		<p>We hope to serve you in the future!</p>
		<br>
		<p>Best regards,<br>IHB Transport Team</p>
	`, deliveryID)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "PRICE_DECLINED", err)
}

// SendDriverOnWayEmail notifies that driver picked up the package
func (s *emailService) SendDriverOnWayEmail(clientEmail, deliveryID string) error {
	subject := fmt.Sprintf("🚚 Driver En Route - Delivery #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2>📦 Great News - Your Package is On the Way!</h2>
		<p>Our driver has successfully picked up your package and is now en route to the delivery location.</p>
		<p><strong>Delivery ID:</strong> %s</p>
		<p><strong>Status:</strong> In Progress</p>
		<p>You can expect delivery soon. We'll send you a confirmation once the package is delivered.</p>
		<p>Thank you for your patience!</p>
		<br>
		<p>Best regards,<br>IHB Transport Team</p>
	`, deliveryID)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "DRIVER_ON_WAY", err)
}

// SendDeliveredEmail sends completion confirmation
func (s *emailService) SendDeliveredEmail(clientEmail, deliveryID string) error {
	subject := fmt.Sprintf("✅ Delivery Completed #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2>🎉 Delivery Complete!</h2>
		<p>Your package has been successfully delivered!</p>
		<p><strong>Delivery ID:</strong> %s</p>
		<p><strong>Status:</strong> Delivered</p>
		<div style="background-color: #d4edda; padding: 15px; margin: 20px 0; border-radius: 5px; border-left: 4px solid #28a745;">
			<p style="margin: 0;"><strong>✅ Delivery confirmed</strong></p>
		</div>
		<p>Thank you for choosing IHB Transport! We hope to serve you again soon.</p>
		<p>If you have any questions or feedback, please don't hesitate to contact us.</p>
		<br>
		<p>Best regards,<br>IHB Transport Team</p>
	`, deliveryID)

	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "DELIVERED", err)
}

// sendEmail sends an email using SMTP
func (s *emailService) sendEmail(to, subject, htmlBody string) error {
	// Check rate limit before attempting to send
	rateLimiter := utils.GetEmailRateLimiter()
	if err := rateLimiter.CanSendEmail(to); err != nil {
		fmt.Printf("⚠️  Rate limit exceeded: %v\n", err)
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	// If SMTP is not enabled, just log to console
	if !s.enabled {
		fmt.Printf("📧 EMAIL (Console): %s\n", subject)
		fmt.Printf("   To: %s\n", to)
		fmt.Printf("   (SMTP not configured - email not sent)\n\n")
		return nil
	}

	// Build the email message
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Setup authentication
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(message))

	if err != nil {
		fmt.Printf("❌ Failed to send email to %s: %v\n", to, err)
		return err
	}

	// Record successful send in rate limiter
	rateLimiter.RecordEmailSent(to)

	fmt.Printf("✅ Email sent successfully to %s: %s\n", to, subject)
	return nil
}

// logEmail creates an email log entry in the database
func (s *emailService) logEmail(deliveryID, recipientEmail, emailType string, emailErr error) error {
	id, err := uuid.Parse(deliveryID)
	if err != nil {
		return fmt.Errorf("invalid delivery ID for email log: %w", err)
	}

	emailLog := models.EmailLog{
		DeliveryID:     id,
		RecipientEmail: recipientEmail,
		EmailType:      emailType,
		Status:         "sent",
	}

	if emailErr != nil {
		emailLog.Status = "failed"
		emailLog.ErrorMessage = emailErr.Error()
	}

	if err := s.emailLogRepo.CreateEmailLog(&emailLog); err != nil {
		fmt.Printf("Warning: failed to log email: %v\n", err)
	}

	return emailErr
}
