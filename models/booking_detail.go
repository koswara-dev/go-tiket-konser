package models

import (
	"github.com/google/uuid"
)

type BookingDetail struct {
	BaseModel
	BookingID        uuid.UUID      `gorm:"type:uuid;not null" json:"booking_id"`
	TicketCategoryID uuid.UUID      `gorm:"type:uuid;not null" json:"ticket_category_id"`
	TicketCategory   TicketCategory `gorm:"foreignKey:TicketCategoryID" json:"ticket_category"`
	Quantity         int            `gorm:"type:integer;not null" json:"quantity"`
	SubTotal         float64        `gorm:"type:numeric(12,2);not null" json:"sub_total"`
}
