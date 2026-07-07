package models

import (
	"time"

	"github.com/google/uuid"
)

// Booking Ticket Concert
type Booking struct {
	BaseModel
	BookingCode string          `gorm:"type:varchar(50);unique;not null" json:"booking_code"`
	CustomerID  uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	Customer    Customer        `gorm:"foreignKey:CustomerID" json:"customer"`
	TotalAmount float64         `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	BookingDate time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"booking_date"`
	Details     []BookingDetail `gorm:"foreignKey:BookingID" json:"details"`
}
