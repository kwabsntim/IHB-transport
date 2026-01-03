package repository

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type deliveryRepository struct {
	db *gorm.DB
}

func NewDeliverRepository() DelivveryInterface {
	return &deliveryRepository{
		db: database.DB,
	}
}

// create a delivery
func (r *deliveryRepository) CreateDelivery(delivery *models.DeliveryRequest) (*models.DeliveryRequest, error) {
	err := r.db.Create(delivery)
	if err != nil {
		log.Fatalf("Could not create delivery")
	}
	return delivery, nil
}

// get all the deliveries

func (r *deliveryRepository) FindAllDeliveries() ([]models.DeliveryRequest, error) {
	var deliveries []models.DeliveryRequest
	//loading the info into the deliveries slice
	err := r.db.
		Preload("StatusLogs", func(db *gorm.DB) *gorm.DB {
			return db.Order("changed_at DESC")
		}).
		Preload("EmailLogs", func(db *gorm.DB) *gorm.DB {
			return db.Order("sent_at DESC")
		}).
		Order("created_at DESC").
		Find(&deliveries).Error
	return deliveries, err
}

// find delivery by ID
func (r *deliveryRepository) FindByID(id uuid.UUID) (*models.DeliveryRequest, error) {
	var delivery models.DeliveryRequest
	err := r.db.
		Preload("StatusLogs", func(db *gorm.DB) *gorm.DB {
			return db.Order("changed_at DESC")
		}).
		Preload("EmailLogs", func(db *gorm.DB) *gorm.DB {
			return db.Order("sent_at DESC")
		}).
		First(&delivery, "id = ?", id).Error

	if err != nil {
		return nil, err // Not found or error
	}
	return &delivery, nil // Found it!
}

// find delivery by status
func (r *deliveryRepository) FindByStatus(status string) ([]models.DeliveryRequest, error) {
	var deliveries []models.DeliveryRequest
	err := r.db.
		Where("status = ?", status).
		Preload("StatusLogs").
		Preload("EmailLogs").
		Order("created_at DESC").
		Find(&deliveries).Error

	return deliveries, err
}

// finds the delivery using its email
func (r *deliveryRepository) FindByEmail(email string) ([]models.DeliveryRequest, error) {
	var deliveries []models.DeliveryRequest
	err := r.db.
		Where("client_email = ?", email).
		Preload("StatusLogs").
		Preload("EmailLogs").
		Order("created_at DESC").
		Find(&deliveries).Error

	return deliveries, err
}

// update the contents of a delivery
func (r *deliveryRepository) UpdateDelivery(delivery *models.DeliveryRequest) error {
	err := r.db.Save(delivery)
	if err != nil {
		log.Fatalf("Could not update delivery")
	}
	return nil
}

// updates the status of a delivery
func (r *deliveryRepository) UpdateStatus(id uuid.UUID, oldStatus, newStatus string) error {
	return r.db.Model(&models.DeliveryRequest{}).
		Where("id = ? AND status = ?", id, oldStatus).
		Update("status", newStatus).Error
}

// updates the price of a delivery
func (r *deliveryRepository) UpdatePrice(id uuid.UUID, price float64) error {
	return r.db.Model(&models.DeliveryRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"price":  price,
			"status": models.StatusPriced,
		}).Error
}

// deletes a delivery
func (r *deliveryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.DeliveryRequest{}, "id = ?", id).Error
}

// counts the total deliveries
func (r *deliveryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.DeliveryRequest{}).Count(&count).Error
	return count, err
}

// counts the number of deliveries by status
func (r *deliveryRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.DeliveryRequest{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
