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

func NewDeliveryRepository() DeliveryInterface {
	return &deliveryRepository{
		db: database.DB,
	}
}

// create a delivery
func (r *deliveryRepository) CreateDelivery(delivery *models.DeliveryRequest) error {
	if err := r.db.Create(delivery).Error; err != nil {
		log.Fatalf("Could not create delivery")
		return err
	}
	return nil
}

// ge
// creates an instant quote
func (r *deliveryRepository) GetInstantQuote(InstantQuote *models.InstantQuote) error {
	if err := r.db.Create(InstantQuote).Error; err != nil {
		log.Printf("Error creating instant quote: %v", err)
		return err
	}
	return nil
}

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
	if err := r.db.Save(delivery).Error; err != nil {
		log.Fatalf("Could not update delivery")
		return err
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

// ===== STATUS LOG REPOSITORY =====

type statusLogRepository struct {
	db *gorm.DB
}

func NewStatusLogRepository() StatusLogInterface {
	return &statusLogRepository{
		db: database.DB,
	}
}

// CreateStatusLog creates a new status log entry
func (r *statusLogRepository) CreateStatusLog(log *models.StatusLog) error {
	return r.db.Create(log).Error
}

// FindByDeliveryID gets all status logs for a delivery
func (r *statusLogRepository) FindByDeliveryID(deliveryID uuid.UUID) ([]models.StatusLog, error) {
	var logs []models.StatusLog
	err := r.db.
		Where("delivery_id = ?", deliveryID).
		Order("changed_at DESC").
		Find(&logs).Error
	return logs, err
}

// FindLatestByDeliveryID gets the most recent status log for a delivery
func (r *statusLogRepository) FindLatestByDeliveryID(deliveryID uuid.UUID) (*models.StatusLog, error) {
	var log models.StatusLog
	err := r.db.
		Where("delivery_id = ?", deliveryID).
		Order("changed_at DESC").
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// ===== EMAIL LOG REPOSITORY =====

type emailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository() EmailLogInterface {
	return &emailLogRepository{
		db: database.DB,
	}
}

// ===== REVIEWS REPOSITORY =====

type reviewsRepository struct {
	db *gorm.DB
}

func NewReviewsRepository() ReviewsInterface {
	return &reviewsRepository{
		db: database.DB,
	}
}

// CreateReview creates a new review entry
func (r *reviewsRepository) CreateReview(review *models.Reviews) error {
	return r.db.Create(review).Error
}

// FindByID finds a review by UUID
func (r *reviewsRepository) FindByID(id uuid.UUID) (*models.Reviews, error) {
	var rev models.Reviews
	if err := r.db.First(&rev, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}

// CreateEmailLog creates a new email log entry
func (r *emailLogRepository) CreateEmailLog(log *models.EmailLog) error {
	return r.db.Create(log).Error
}

// FindByDeliveryID gets all email logs for a delivery
func (r *emailLogRepository) FindByDeliveryID(deliveryID uuid.UUID) ([]models.EmailLog, error) {
	var logs []models.EmailLog
	err := r.db.
		Where("delivery_id = ?", deliveryID).
		Order("sent_at DESC").
		Find(&logs).Error
	return logs, err
}

// UpdateEmailStatus updates the status of an email log (e.g., from "pending" to "sent")
func (r *emailLogRepository) UpdateEmailStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.EmailLog{}).
		Where("id = ?", id).
		Update("status", status).Error
}
