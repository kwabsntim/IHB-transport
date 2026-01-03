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
