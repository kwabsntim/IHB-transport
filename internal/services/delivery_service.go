package services

import (
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
	"ihb-transport/utils"
	"time"

	"github.com/google/uuid"
)

type deliveryService struct {
	deliveryRepo  repository.DeliveryInterface
	statusLogRepo repository.StatusLogInterface
	emailLogRepo  repository.EmailLogInterface
	emailService  EmailServiceInterface
}

func NewDeliveryService(
	deliveryRepo repository.DeliveryInterface,
	statusLogRepo repository.StatusLogInterface,
	emailLogRepo repository.EmailLogInterface,
	emailService EmailServiceInterface,
) DeliveryServiceInterface {
	return &deliveryService{
		deliveryRepo:  deliveryRepo,
		statusLogRepo: statusLogRepo,
		emailLogRepo:  emailLogRepo,
		emailService:  emailService,
	}
}

//delivery service functions

// validation functions for creating a delivery request
func (s *deliveryService) validateCreateInput(delivery *models.DeliveryRequest) error {
	// Validate required fields
	if err := utils.ValidateRequired("client name", delivery.ClientName); err != nil {
		return err
	}

	if err := utils.ValidateRequired("client email", delivery.ClientEmail); err != nil {
		return err
	}

	// Validate pickup address fields
	if err := utils.ValidateRequired("pickup street", delivery.PickupStreet); err != nil {
		return err
	}
	if err := utils.ValidateRequired("pickup city", delivery.PickupCity); err != nil {
		return err
	}
	if err := utils.ValidateRequired("pickup post code", delivery.PickupPostCode); err != nil {
		return err
	}
	if err := utils.ValidateRequired("pickup country", delivery.PickupCountry); err != nil {
		return err
	}

	// Validate dropoff address fields
	if err := utils.ValidateRequired("dropoff street", delivery.DropoffStreet); err != nil {
		return err
	}
	if err := utils.ValidateRequired("dropoff city", delivery.DropoffCity); err != nil {
		return err
	}
	if err := utils.ValidateRequired("dropoff post code", delivery.DropoffPostCode); err != nil {
		return err
	}
	if err := utils.ValidateRequired("dropoff country", delivery.DropoffCountry); err != nil {
		return err
	}

	// Validate email format
	if err := utils.ValidateEmail(delivery.ClientEmail); err != nil {
		return err
	}

	// Validate items description
	if err := utils.ValidateRequired("items", delivery.Items); err != nil {
		return err
	}

	// Validate service
	if err := utils.ValidateRequired("service", delivery.Service); err != nil {
		return err
	}

	// Validate description length
	if err := utils.ValidateMaxLength("item description", delivery.ItemDescription, 500); err != nil {
		return err
	}

	return nil
}

// becaus you are already using pointer in &delivery use delivery
func (s *deliveryService) CreateDeliveryRequest(delivery *models.DeliveryRequest) error {
	//validating the delivery input
	if err := s.validateCreateInput(delivery); err != nil {
		return err
	}
	//modifying the status of the delivery request
	delivery.Status = models.StatusRequested
	delivery.ID = uuid.Nil
	//creating the delivery request
	if err := s.deliveryRepo.CreateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to create a delivery request: %w", err)
	}
	//logging the status of the delivery request
	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  "",
		NewStatus:  models.StatusRequested,
		ChangedBy:  "system",
	}
	//saving the status log
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	// Send confirmation email to client
	pickupAddr := fmt.Sprintf("%s, %s, %s, %s", delivery.PickupStreet, delivery.PickupCity, delivery.PickupPostCode, delivery.PickupCountry)
	dropoffAddr := fmt.Sprintf("%s, %s, %s, %s", delivery.DropoffStreet, delivery.DropoffCity, delivery.DropoffPostCode, delivery.DropoffCountry)
	pickupDateStr := "Not specified"
	if delivery.PickupDate != nil {
		pickupDateStr = delivery.PickupDate.Format("Monday, January 2, 2006")
	}

	if err := s.emailService.SendRequestReceivedEmail(
		delivery.ClientEmail,
		delivery.ClientName,
		delivery.ID.String(),
		pickupAddr,
		dropoffAddr,
		delivery.Service,
		pickupDateStr,
	); err != nil {
		fmt.Printf("Warning: failed to send request received email: %v\n", err)
	}

	return nil
}

// set the delivery price for a delivery request
func (s *deliveryService) SetDeliveryPrice(id string, price float64) error {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID format: %w", err)
	}
	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("failed to find delivery by ID: %w", err)
	}
	if delivery.Status != models.StatusRequested {
		return fmt.Errorf("cannot set price for delivery with status %s", delivery.Status)
	}
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	oldstatus := delivery.Status
	delivery.Price = price
	delivery.Status = models.StatusPriced

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	// Log status change
	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldstatus,
		NewStatus:  models.StatusPriced,
		ChangedBy:  "admin",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	// Send price email
	if err := s.emailService.SendPriceEmail(delivery.ClientEmail, price, delivery.ID.String()); err != nil {
		fmt.Printf("Warning: failed to send price email: %v\n", err)
	}

	return nil
}

// DeclineDeliveryPrice - Client declines the quoted price
func (s *deliveryService) DeclineDeliveryPrice(id string, reason string) error {
	// Step 1: Parse UUID
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID: %w", err)
	}

	// Step 2: Find delivery
	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	// Step 3: Validate status (must be PRICED)
	if delivery.Status != models.StatusPriced {
		// Provide user-friendly error messages based on current status
		switch delivery.Status {
		case models.StatusAccepted:
			return fmt.Errorf("this quote has already been accepted and cannot be declined")
		case models.StatusDeclined:
			return fmt.Errorf("this quote has already been declined")
		case models.StatusInProgress:
			return fmt.Errorf("this delivery has already been picked up and cannot be declined")
		case models.StatusDelivered:
			return fmt.Errorf("this delivery has already been completed")
		case models.StatusPending:
			return fmt.Errorf("no price has been set for this delivery yet")
		default:
			return fmt.Errorf("cannot decline delivery in current status: %s", delivery.Status)
		}
	}

	// Step 4: Validate reason is provided and has reasonable length
	if reason == "" {
		return fmt.Errorf("decline reason is required")
	}
	// Limit reason length to prevent abuse (database field is text, but we limit to 1000 chars)
	if len(reason) > 1000 {
		return fmt.Errorf("decline reason cannot exceed 1000 characters")
	}

	// Step 5: Update delivery with decline info
	oldStatus := delivery.Status
	delivery.Status = models.StatusDeclined
	delivery.DeclineReason = reason
	now := time.Now()
	delivery.DeclinedAt = &now

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	// Step 6: Log status change
	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldStatus,
		NewStatus:  models.StatusDeclined,
		ChangedBy:  "client",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	// Step 7: Send email notification about decline
	if err := s.emailService.SendDeclinedEmail(delivery.ClientEmail, delivery.ID.String(), reason); err != nil {
		fmt.Printf("Warning: failed to send decline email: %v\n", err)
	}

	return nil
}

// AcceptDeliveryPrice - Client accepts the quoted price
func (s *deliveryService) AcceptDeliveryPrice(id string) error {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID: %w", err)
	}

	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	// Validate status (must be PRICED)
	if delivery.Status != models.StatusPriced {
		// Provide user-friendly error messages based on current status
		switch delivery.Status {
		case models.StatusAccepted:
			return fmt.Errorf("this quote has already been accepted")
		case models.StatusDeclined:
			return fmt.Errorf("this quote has already been declined and cannot be accepted")
		case models.StatusInProgress:
			return fmt.Errorf("this delivery has already been picked up")
		case models.StatusDelivered:
			return fmt.Errorf("this delivery has already been completed")
		case models.StatusPending:
			return fmt.Errorf("no price has been set for this delivery yet")
		default:
			return fmt.Errorf("cannot accept delivery in current status: %s", delivery.Status)
		}
	}

	oldStatus := delivery.Status
	delivery.Status = models.StatusAccepted

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldStatus,
		NewStatus:  models.StatusAccepted,
		ChangedBy:  "client",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	if err := s.emailService.SendAcceptedEmail(delivery.ClientEmail, delivery.ID.String()); err != nil {
		fmt.Printf("Warning: failed to send accepted email: %v\n", err)
	}

	return nil
}

// AdminAcceptDeliveryPrice - Admin accepts the delivery on behalf of client
func (s *deliveryService) AdminAcceptDeliveryPrice(id string) error {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID: %w", err)
	}

	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	// Validate status (must be PRICED)
	if delivery.Status != models.StatusPriced {
		switch delivery.Status {
		case models.StatusAccepted:
			return fmt.Errorf("this quote has already been accepted")
		case models.StatusDeclined:
			return fmt.Errorf("this quote has already been declined and cannot be accepted")
		case models.StatusInProgress:
			return fmt.Errorf("this delivery has already been picked up")
		case models.StatusDelivered:
			return fmt.Errorf("this delivery has already been completed")
		case models.StatusPending:
			return fmt.Errorf("no price has been set for this delivery yet")
		default:
			return fmt.Errorf("cannot accept delivery in current status: %s", delivery.Status)
		}
	}

	oldStatus := delivery.Status
	delivery.Status = models.StatusAccepted

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldStatus,
		NewStatus:  models.StatusAccepted,
		ChangedBy:  "admin",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	if err := s.emailService.SendAcceptedEmail(delivery.ClientEmail, delivery.ID.String()); err != nil {
		fmt.Printf("Warning: failed to send accepted email: %v\n", err)
	}

	return nil
}

// MarkAsPickedUp - Driver marks delivery as picked up
func (s *deliveryService) MarkAsPickedUp(id string) error {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID: %w", err)
	}

	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	if delivery.Status != models.StatusAccepted {
		return fmt.Errorf("can only pick up ACCEPTED deliveries, current status: %s", delivery.Status)
	}

	oldStatus := delivery.Status
	delivery.Status = models.StatusInProgress

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldStatus,
		NewStatus:  models.StatusInProgress,
		ChangedBy:  "driver",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	if err := s.emailService.SendDriverOnWayEmail(delivery.ClientEmail, delivery.ID.String()); err != nil {
		fmt.Printf("Warning: failed to send driver on way email: %v\n", err)
	}

	return nil
}

// MarkAsDelivered - Driver marks delivery as completed
func (s *deliveryService) MarkAsDelivered(id string) error {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid delivery ID: %w", err)
	}

	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	if delivery.Status != models.StatusInProgress {
		return fmt.Errorf("can only complete IN_PROGRESS deliveries, current status: %s", delivery.Status)
	}

	oldStatus := delivery.Status
	delivery.Status = models.StatusDelivered

	if err := s.deliveryRepo.UpdateDelivery(delivery); err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	statusLog := models.StatusLog{
		DeliveryID: delivery.ID,
		OldStatus:  oldStatus,
		NewStatus:  models.StatusDelivered,
		ChangedBy:  "driver",
	}
	if err := s.statusLogRepo.CreateStatusLog(&statusLog); err != nil {
		fmt.Printf("Warning: failed to create status log: %v\n", err)
	}

	if err := s.emailService.SendDeliveredEmail(delivery.ClientEmail, delivery.ID.String()); err != nil {
		fmt.Printf("Warning: failed to send delivered email: %v\n", err)
	}

	return nil
}

// GetAllDeliveries - Get all deliveries
func (s *deliveryService) GetAllDeliveries() ([]models.DeliveryRequest, error) {
	deliveries, err := s.deliveryRepo.FindAllDeliveries()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}
	return deliveries, nil
}

// GetDeliveryByID - Get delivery by ID
func (s *deliveryService) GetDeliveryByID(id string) (*models.DeliveryRequest, error) {
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid delivery ID: %w", err)
	}

	delivery, err := s.deliveryRepo.FindByID(deliveryID)
	if err != nil {
		return nil, fmt.Errorf("delivery not found: %w", err)
	}

	return delivery, nil
}

// GetDeliveriesByEmail - Get all deliveries for a client
func (s *deliveryService) GetDeliveriesByEmail(email string) ([]models.DeliveryRequest, error) {
	if err := utils.ValidateEmail(email); err != nil {
		return nil, err
	}

	deliveries, err := s.deliveryRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}

	return deliveries, nil
}

// GetDeliveriesByStatus - Filter deliveries by status
func (s *deliveryService) GetDeliveriesByStatus(status string) ([]models.DeliveryRequest, error) {
	validStatuses := []string{
		models.StatusRequested,
		models.StatusPriced,
		models.StatusDeclined,
		models.StatusAccepted,
		models.StatusInProgress,
		models.StatusDelivered,
	}

	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	deliveries, err := s.deliveryRepo.FindByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %w", err)
	}

	return deliveries, nil
}
func (s *deliveryService) GetInstantQuote(instantQuote *models.InstantQuote) error {
	if err := utils.ValidateRequired("pickup point", instantQuote.PickupPoint); err != nil {
		return err
	}
	if err := utils.ValidateRequired("delivery address", instantQuote.DeliveryAddress); err != nil {
		return err
	}
	if err := utils.ValidateRequired("weight", instantQuote.Weight); err != nil {
		return err
	}
	if err := utils.ValidateEmail(instantQuote.ClientEmail); err != nil {
		return err
	}
	// Here you would implement your logic to calculate the instant quote based on the provided details.
	// For now, we will just return a success message.
	if err := s.deliveryRepo.GetInstantQuote(instantQuote); err != nil {
		return err
	}
	if err := s.emailService.SendInstantQuoteEmail(instantQuote.ClientEmail, instantQuote.ID.String(), instantQuote.PickupPoint, instantQuote.DeliveryAddress, instantQuote.Weight); err != nil {
		fmt.Printf("Warning: failed to send instant quote email: %v\n", err)
	}
	return nil

}
