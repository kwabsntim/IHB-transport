package services

import (
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"

	"github.com/google/uuid"
)

type emailService struct {
	emailLogRepo repository.EmailLogInterface
}

// NewEmailService creates a new email service
func NewEmailService(emailLogRepo repository.EmailLogInterface) EmailServiceInterface {
	return &emailService{
		emailLogRepo: emailLogRepo,
	}
}

// SendRequestReceivedEmail sends confirmation email when request is created
func (s *emailService) SendRequestReceivedEmail(clientEmail, clientName, deliveryID string) error {
	fmt.Printf("📧 EMAIL: REQUEST_RECEIVED\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Delivery Request Received #%s\n", deliveryID)
	fmt.Printf("   Body: Hello %s, we received your delivery request.\n\n", clientName)

	return s.logEmail(deliveryID, clientEmail, "REQUEST_RECEIVED", nil)
}

// SendPriceEmail sends email with quoted price
func (s *emailService) SendPriceEmail(clientEmail string, price float64, deliveryID string) error {
	fmt.Printf("📧 EMAIL: PRICE_SENT\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Price Quote for Delivery #%s\n", deliveryID)
	fmt.Printf("   Body: Your delivery quote is $%.2f\n\n", price)

	return s.logEmail(deliveryID, clientEmail, "PRICE_SENT", nil)
}

// SendAcceptedEmail notifies that client accepted the price
func (s *emailService) SendAcceptedEmail(clientEmail, deliveryID string) error {
	fmt.Printf("📧 EMAIL: ACCEPTED\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Delivery Confirmed #%s\n", deliveryID)
	fmt.Printf("   Body: Thank you for accepting! Your delivery is confirmed.\n\n")

	return s.logEmail(deliveryID, clientEmail, "ACCEPTED", nil)
}

// SendDeclinedEmail sends notification when client declines the price
func (s *emailService) SendDeclinedEmail(clientEmail, deliveryID, reason string) error {
	fmt.Printf("📧 EMAIL: PRICE_DECLINED\n")
	fmt.Printf("   To: admin@ihbtransport.com\n")
	fmt.Printf("   Subject: Price Declined - Delivery #%s\n", deliveryID)
	fmt.Printf("   Body: Client %s declined the price.\n", clientEmail)
	fmt.Printf("   Reason: %s\n\n", reason)

	fmt.Printf("📧 EMAIL: DECLINE_CONFIRMATION\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Delivery Request Declined #%s\n", deliveryID)
	fmt.Printf("   Body: We understand the price wasn't suitable. We may contact you with a better offer.\n\n")

	return s.logEmail(deliveryID, clientEmail, "PRICE_DECLINED", nil)
}

// SendDriverOnWayEmail notifies that driver picked up the package
func (s *emailService) SendDriverOnWayEmail(clientEmail, deliveryID string) error {
	fmt.Printf("📧 EMAIL: DRIVER_ON_WAY\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Driver En Route - Delivery #%s\n", deliveryID)
	fmt.Printf("   Body: Great news! Our driver has picked up your package and is on the way.\n\n")

	return s.logEmail(deliveryID, clientEmail, "DRIVER_ON_WAY", nil)
}

// SendDeliveredEmail sends completion confirmation
func (s *emailService) SendDeliveredEmail(clientEmail, deliveryID string) error {
	fmt.Printf("📧 EMAIL: DELIVERED\n")
	fmt.Printf("   To: %s\n", clientEmail)
	fmt.Printf("   Subject: Delivery Completed #%s\n", deliveryID)
	fmt.Printf("   Body: Your package has been successfully delivered. Thank you!\n\n")

	return s.logEmail(deliveryID, clientEmail, "DELIVERED", nil)
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
