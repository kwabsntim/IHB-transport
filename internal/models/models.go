package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Email        string    `gorm:"type:varchar(120);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"password_hash"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// delivery request model
type DeliveryRequest struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ClientName      string      `gorm:"type:varchar(120)" json:"client_name"`
	ClientEmail     string      `gorm:"type:varchar(120);not null" json:"client_email"`
	PickupAddress   string      `gorm:"type:varchar(255);not null" json:"pickup_address"`
	DropoffAddress  string      `gorm:"type:varchar(255);not null" json:"dropoff_address"`
	ItemDescription string      `gorm:"type:varchar(255)" json:"item_description"`
	Weight          float64     `gorm:"type:decimal(10,2)" json:"weight"`
	Status          string      `gorm:"type:varchar(20);default:'REQUESTED'" json:"status"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	Price           float64     `gorm:"type:decimal(10,2);default:0.0" json:"price"`
	StatusLogs      []StatusLog `gorm:"foreignKey:DeliveryID" json:"status_logs,omitempty"`
	EmailLogs       []EmailLog  `gorm:"foreignKey:DeliveryID" json:"email_logs,omitempty"`
}

// status log model
type StatusLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DeliveryID uuid.UUID `gorm:"type:uuid;not null;index" json:"delivery_id"`
	OldStatus  string    `gorm:"type:varchar(20)" json:"old_status"`
	NewStatus  string    `gorm:"type:varchar(20)" json:"new_status"`
	ChangedBy  string    `gorm:"type:varchar(50)" json:"changed_by"` // "driver" or "system"
	ChangedAt  time.Time `gorm:"autoCreateTime" json:"changed_at"`
}

// email log model
type EmailLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DeliveryID     uuid.UUID `gorm:"type:uuid;not null;index" json:"delivery_id"`
	RecipientEmail string    `gorm:"type:varchar(120);not null" json:"recipient_email"`
	EmailType      string    `gorm:"type:varchar(50)" json:"email_type"` // REQUEST_RECEIVED, PRICE_SENT, DRIVER_ON_WAY, DELIVERED
	Status         string    `gorm:"type:varchar(20)" json:"status"`     // sent, failed
	SentAt         time.Time `gorm:"autoCreateTime" json:"sent_at"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
}

// constants for the status
const (
	StatusRequested  = "REQUESTED"
	StatusPriced     = "PRICED"
	StatusAccepted   = "ACCEPTED"
	StatusInProgress = "IN_PROGRESS"
	StatusDelivered  = "DELIVERED"
)
