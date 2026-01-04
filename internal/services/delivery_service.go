package services

import (
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
	"ihb-transport/utils"

	"github.com/google/uuid"
)

type deliveryService struct {
	deliveryRepo  repository.DeliveryInterface
	statusLogRepo repository.StatusLogInterface
	emailLogRepo  repository.EmailLogInterface
}

func NewDeliveryService(
	deliveryRepo repository.DeliveryInterface,
	statusLogRepo repository.StatusLogInterface,
	emailLogRepo repository.EmailLogInterface,
) DeliveryServiceInterface {
	return &deliveryService{
		deliveryRepo:  deliveryRepo,
		statusLogRepo: statusLogRepo,
		emailLogRepo:  emailLogRepo,
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

	if err := utils.ValidateRequired("pickup address", delivery.PickupAddress); err != nil {
		return err
	}

	if err := utils.ValidateRequired("dropoff address", delivery.DropoffAddress); err != nil {
		return err
	}

	// Validate email format
	if err := utils.ValidateEmail(delivery.ClientEmail); err != nil {
		return err
	}

	// Validate weight is not negative
	if err := utils.ValidatePositiveNumber("weight", delivery.Weight); err != nil {
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

}
