package repository

import (
	"ihb-transport/internal/models"

	"github.com/google/uuid"
)

type DelivveryInterface interface {
	// Create operations
	Create(delivery *models.DeliveryRequest) error

	// Read operations
	FindAll() ([]models.DeliveryRequest, error)
	FindByID(id uuid.UUID) (*models.DeliveryRequest, error)
	FindByStatus(status string) ([]models.DeliveryRequest, error)
	FindByEmail(email string) ([]models.DeliveryRequest, error)

	// Update operations
	Update(delivery *models.DeliveryRequest) error
	UpdateStatus(id uuid.UUID, oldStatus, newStatus string) error
	UpdatePrice(id uuid.UUID, price float64) error

	// Delete operations
	Delete(id uuid.UUID) error

	// Counting operations
	Count() (int64, error)
	CountByStatus(status string) (int64, error)
}
