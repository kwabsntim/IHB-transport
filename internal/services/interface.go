package services

import "ihb-transport/internal/models"

// DeliveryServiceInterface defines delivery business logic operations
type DeliveryServiceInterface interface {
	CreateDeliveryRequest(delivery *models.DeliveryRequest) error
	SetDeliveryPrice(id string, price float64) error
	AcceptDeliveryPrice(id string) error
	DeclineDeliveryPrice(id string, reason string) error
	MarkAsPickedUp(id string) error
	MarkAsDelivered(id string) error
	GetAllDeliveries() ([]models.DeliveryRequest, error)
	GetDeliveryByID(id string) (*models.DeliveryRequest, error)
	GetDeliveriesByEmail(email string) ([]models.DeliveryRequest, error)
	GetDeliveriesByStatus(status string) ([]models.DeliveryRequest, error)
	GetInstantQuote(instantQuote *models.InstantQuote) error
}

// EmailServiceInterface defines email operations
type EmailServiceInterface interface {
	SendRequestReceivedEmail(clientEmail, clientName, deliveryID, pickupAddress, dropoffAddress, service, pickupDate string) error
	SendPriceEmail(clientEmail string, price float64, deliveryID string) error
	SendAcceptedEmail(clientEmail, deliveryID string) error
	SendDeclinedEmail(clientEmail, deliveryID, reason string) error
	SendDriverOnWayEmail(clientEmail, deliveryID string) error
	SendDeliveredEmail(clientEmail, deliveryID string) error
	SendInstantQuoteEmail(clientEmail, quoteID, pickupPoint, deliveryAddress, weight string) error
}

// ReviewServiceInterface defines review operations
type ReviewServiceInterface interface {
	CreateReview(review *models.Reviews) error
	GetReviewByID(id string) (*models.Reviews, error)
}
