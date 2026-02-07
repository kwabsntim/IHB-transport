package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
	"ihb-transport/internal/utils"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

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
	// Resend (HTTP API) support
	resendAPIKey  string
	resendEnabled bool
	// Frontend base URL (for links that should open the SPA)
	frontendURL string
}

// NewEmailService creates a new email service
func NewEmailService(emailLogRepo repository.EmailLogInterface) EmailServiceInterface {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	// Resend settings
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	resendFromEmail := os.Getenv("RESEND_FROM_EMAIL")
	resendFromName := os.Getenv("RESEND_FROM_NAME")

	// Check if SMTP is enabled
	enabled := smtpHost != "" && smtpPort != "" && smtpUser != "" && smtpPassword != ""

	if !enabled {
		fmt.Println("⚠️  SMTP not configured - emails will be logged to console only")
		fmt.Println("   To enable SMTP, set: SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD")
	} else {
		fmt.Printf("✅ SMTP enabled - using %s:%s\n", smtpHost, smtpPort)
	}

	// Prefer explicit RESEND_FROM_EMAIL if provided, else fallback to SMTP_FROM_EMAIL
	if resendFromEmail != "" {
		fromEmail = resendFromEmail
	}
	if resendFromName != "" {
		fromName = resendFromName
	}
	if fromEmail == "" {
		fromEmail = "noreply@ihbtransport.com"
	}
	if fromName == "" {
		fromName = "IHB Transport"
	}

	// Get base URL for backend links used in email templates
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	// Frontend URL (optional) — used when you want the email to point to your SPA
	frontendURL := os.Getenv("FRONTEND_BASE_URL")

	// Resend enabled when API key is present
	resendEnabled := resendAPIKey != ""

	return &emailService{
		emailLogRepo:  emailLogRepo,
		smtpHost:      smtpHost,
		smtpPort:      smtpPort,
		smtpUser:      smtpUser,
		smtpPassword:  smtpPassword,
		fromEmail:     fromEmail,
		fromName:      fromName,
		baseURL:       baseURL,
		frontendURL:   frontendURL,
		enabled:       enabled,
		resendAPIKey:  resendAPIKey,
		resendEnabled: resendEnabled,
	}
}

// resendPayload is the request body we send to Resend API
type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// sendWithResend sends the email using Resend's HTTP API
func (s *emailService) sendWithResend(to []string, subject, html string) error {
	if s.resendAPIKey == "" {
		return fmt.Errorf("resend api key not configured")
	}

	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	payload := resendPayload{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal resend payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create resend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.resendAPIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("resend request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("resend returned status %d", resp.StatusCode)
	}

	return nil
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
	//this sends the email to the client and logs it
	err := s.sendEmail(clientEmail, subject, body)
	return s.logEmail(deliveryID, clientEmail, "REQUEST_RECEIVED", err)
}
func (s *emailService) SendInstantQuoteEmail(clientEmail, quoteID, pickupPoint, deliveryAddress, weight string) error {
	// Require recipient email to send instant quote
	if strings.TrimSpace(clientEmail) == "" {
		return fmt.Errorf("recipient email is required to send instant quote")
	}

	subject := fmt.Sprintf("Your Instant Quote is Ready - ID: %s", quoteID)

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		  <meta charset="utf-8" />
		  <style>
		    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		    .container { max-width: 600px; margin: 0 auto; padding: 20px; }
		    .header { background-color: #667eea; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
		    .content { background-color: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
		    .info-box { background-color: white; padding: 15px; margin: 10px 0; border-left: 4px solid #667eea; }
		    .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
		    ul { list-style: none; padding: 0; }
		    li { padding: 8px 0; border-bottom: 1px solid #eee; }
		    li:last-child { border-bottom: none; }
		  </style>
		</head>
		<body>
		  <div class="container">
		    <div class="header">
		      <h1 style="margin: 0;">🚀 Instant Quote Ready!</h1>
		    </div>
		    <div class="content">
		      <p>Dear Customer,</p>
		      <p>Thank you for using our Instant Quote feature! Here are the details of your quote:</p>
		      
		      <div class="info-box">
		        <ul>
		          <li><strong>Quote ID:</strong> %s</li>
		          <li><strong>Pickup Point:</strong> %s</li>
		          <li><strong>Delivery Address:</strong> %s</li>
		          <li><strong>Weight:</strong> %s</li>
		        </ul>
		      </div>
		      
		      <p>If you have any questions or would like to proceed with this quote, please contact our support team.</p>
		    </div>
		    <div class="footer">
		      <p>Best regards,<br><strong>IHB Transport Team</strong></p>
		      <p>IHB Transport APS - Reliable Delivery Services</p>
		    </div>
		  </div>
		</body>
		</html>
	`, quoteID, pickupPoint, deliveryAddress, weight)
	//logging the email

	err := s.sendEmail(clientEmail, subject, body)
	if err != nil {
		return err
	}
	return s.logEmail(quoteID, clientEmail, "QUOTE_RECEIVED", err)
}

// SendPriceEmail sends email with quoted price
func (s *emailService) SendPriceEmail(clientEmail string, price float64, deliveryID string) error {
	subject := fmt.Sprintf("Price Quote for Delivery #%s", deliveryID)
	// Build frontend and backend accept/decline links. Prefer frontend if configured.
	backendAccept := fmt.Sprintf("%s/api/public/deliveries/%s/accept", s.baseURL, deliveryID)
	backendDecline := fmt.Sprintf("%s/api/public/deliveries/%s/decline", s.baseURL, deliveryID)
	acceptLink := backendAccept
	declineLink := backendDecline
	if strings.TrimSpace(s.frontendURL) != "" {
		acceptLink = fmt.Sprintf("%s/deliveries/%s/accept", strings.TrimRight(s.frontendURL, "/"), deliveryID)
		declineLink = fmt.Sprintf("%s/deliveries/%s/decline", strings.TrimRight(s.frontendURL, "/"), deliveryID)
	}

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
						<a href="%s" class="button accept-btn">✅ Accept Quote</a>
						<a href="%s" class="button decline-btn">❌ Decline Quote</a>
					</div>
					
					<p style="font-size: 12px; color: #666;">This quote is valid for 48 hours.</p>
				</div>
				<div class="footer">
					<p>IHB Transport - Reliable Delivery Services</p>
				</div>
			</div>
		</body>
		</html>
	`, deliveryID, price, acceptLink, declineLink)

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
	subject := fmt.Sprintf(" Driver En Route - Delivery #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2> Great News - Your Package is On the Way!</h2>
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
	subject := fmt.Sprintf(" Delivery Completed #%s", deliveryID)
	body := fmt.Sprintf(`
		<h2> Delivery Complete!</h2>
		<p>Your package has been successfully delivered!</p>
		<p><strong>Delivery ID:</strong> %s</p>
		<p><strong>Status:</strong> Delivered</p>
		<div style="background-color: #d4edda; padding: 15px; margin: 20px 0; border-radius: 5px; border-left: 4px solid #28a745;">
			<p style="margin: 0;"><strong> Delivery confirmed</strong></p>
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
		fmt.Printf("  Rate limit exceeded: %v\n", err)
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	// If Resend is configured prefer it (HTTP API) — avoids SMTP egress issues
	if s.resendEnabled {
		if err := s.sendWithResend([]string{to}, subject, htmlBody); err != nil {
			fmt.Printf("❌ Resend send failed to %s: %v\n", to, err)
			return err
		}
		// Record successful send in rate limiter
		rateLimiter.RecordEmailSent(to)
		fmt.Printf("✅ Email sent via Resend to %s: %s\n", to, subject)
		return nil
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
