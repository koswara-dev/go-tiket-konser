package models

import (
	"time"

	"gorm.io/gorm"
)

// Booking Ticket Concert
type Booking struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	BookingCode string          `gorm:"type:varchar(50);unique;not null" json:"booking_code"`
	CustomerID  uint            `gorm:"not null" json:"customer_id"`
	Customer    Customer        `gorm:"foreignKey:CustomerID" json:"customer"`
	TotalAmount float64         `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	BookingDate time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"booking_date"`
	Details     []BookingDetail `gorm:"foreignKey:BookingID" json:"details"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}
