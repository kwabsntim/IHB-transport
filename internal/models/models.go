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
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ClientName  string    `gorm:"type:varchar(120)" json:"client_name"`
	ClientEmail string    `gorm:"type:varchar(120);not null" json:"client_email"`

	// Pickup Address Details
	PickupStreet   string `gorm:"type:varchar(255)" json:"pickup_street"`
	PickupCity     string `gorm:"type:varchar(100)" json:"pickup_city"`
	PickupPostCode string `gorm:"type:varchar(20)" json:"pickup_post_code"`
	PickupCountry  string `gorm:"type:varchar(100)" json:"pickup_country"`

	// Dropoff Address Details
	DropoffStreet   string `gorm:"type:varchar(255)" json:"dropoff_street"`
	DropoffCity     string `gorm:"type:varchar(100)" json:"dropoff_city"`
	DropoffPostCode string `gorm:"type:varchar(20)" json:"dropoff_post_code"`
	DropoffCountry  string `gorm:"type:varchar(100)" json:"dropoff_country"`

	// Item Details
	ItemDescription string `gorm:"type:varchar(255)" json:"item_description"`
	Items           string `gorm:"type:text" json:"items"`         // Description of items being delivered
	Weight          string `gorm:"type:varchar(50)" json:"weight"` // Weight of items (e.g., "5kg", "10 pounds")

	// Service Details
	Service    string     `gorm:"type:varchar(100)" json:"service"` // Service type selected from dropdown
	PickupDate *time.Time `gorm:"type:date" json:"pickup_date"`     // Preferred pickup date

	// Delivery Details
	Price         float64     `gorm:"type:decimal(10,2);default:0.0" json:"price"`
	PriceOption   string      `gorm:"type:varchar(100)" json:"price_option"`
	Status        string      `gorm:"type:varchar(20);default:'REQUESTED'" json:"status"`
	DeclineReason string      `gorm:"type:text" json:"decline_reason,omitempty"`
	DeclinedAt    *time.Time  `gorm:"type:timestamp" json:"declined_at,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	StatusLogs    []StatusLog `gorm:"foreignKey:DeliveryID" json:"status_logs,omitempty"`
	EmailLogs     []EmailLog  `gorm:"foreignKey:DeliveryID" json:"email_logs,omitempty"`
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
type InstantQuote struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PickupPoint     string     `json:"pickup_point"`
	DeliveryAddress string     `json:"delivery_address"`
	Weight          string     `json:"weight"`
	ClientEmail     string     `json:"client_email"`
	Price           float64    `gorm:"type:decimal(10,2);default:0.0" json:"price"`
	Status          string     `gorm:"type:varchar(20);default:'REQUESTED'" json:"status"`
	DeclineReason   string     `gorm:"type:text" json:"decline_reason,omitempty"`
	DeclinedAt      *time.Time `gorm:"type:timestamp" json:"declined_at,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
type Reviews struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ClientName string    `gorm:"type:varchar(120)" json:"client_name"`
	Content    string    `gorm:"type:text" json:"content"`
}

// constants for the status
const (
	StatusRequested  = "REQUESTED"
	StatusPriced     = "PRICED"
	StatusDeclined   = "DECLINED"
	StatusAccepted   = "ACCEPTED"
	StatusInProgress = "IN_PROGRESS"
	StatusDelivered  = "DELIVERED"
	StatusPending    = "PENDING"
)
