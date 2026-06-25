package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingDetail struct {
	ID               int            `gorm:"primaryKey" json:"id"`
	BookingID        int            `gorm:"not null" json:"booking_id"`
	TicketCategoryID int            `gorm:"not null" json:"ticket_category_id"`
	TicketCategory   TicketCategory `gorm:"foreignKey:TicketCategoryID" json:"ticket_category"`
	Quantity         int            `gorm:"type:integer;not null" json:"quantity"`
	SubTotal         float64        `gorm:"type:numeric(12,2);not null" json:"sub_total"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
