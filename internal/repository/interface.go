package repository

import (
	"ihb-transport/internal/models"

	"github.com/google/uuid"
)

type DeliveryInterface interface {
	// Create operations
	CreateDelivery(delivery *models.DeliveryRequest) error // ← Match implementation

	// Read operations
	FindAllDeliveries() ([]models.DeliveryRequest, error) // ← Match implementation
	FindByID(id uuid.UUID) (*models.DeliveryRequest, error)
	FindByStatus(status string) ([]models.DeliveryRequest, error)
	FindByEmail(email string) ([]models.DeliveryRequest, error)
	GetInstantQuote(InstantQuote *models.InstantQuote) error

	// InstantQuote operations
	FindInstantQuoteByID(id uuid.UUID) (*models.InstantQuote, error)
	FindAllInstantQuotes() ([]models.InstantQuote, error)
	UpdateInstantQuote(quote *models.InstantQuote) error

	// Update operations
	UpdateDelivery(delivery *models.DeliveryRequest) error // ← Match implementation
	UpdateStatus(id uuid.UUID, oldStatus, newStatus string) error
	UpdatePrice(id uuid.UUID, price float64) error

	// Delete operations
	Delete(id uuid.UUID) error

	// Counting operations
	Count() (int64, error)
	CountByStatus(status string) (int64, error)
}

// StatusLogInterface defines operations for status log repository
type StatusLogInterface interface {
	CreateStatusLog(log *models.StatusLog) error
	FindByDeliveryID(deliveryID uuid.UUID) ([]models.StatusLog, error)
	FindLatestByDeliveryID(deliveryID uuid.UUID) (*models.StatusLog, error)
}

// EmailLogInterface defines operations for email log repository
type EmailLogInterface interface {
	CreateEmailLog(log *models.EmailLog) error
	FindByDeliveryID(deliveryID uuid.UUID) ([]models.EmailLog, error)
	UpdateEmailStatus(id uuid.UUID, status string) error
}

// ReviewsInterface defines operations for reviews
type ReviewsInterface interface {
	CreateReview(review *models.Reviews) error
	FindByID(id uuid.UUID) (*models.Reviews, error)
	FindAll() ([]models.Reviews, error)
	Delete(id uuid.UUID) error
}
