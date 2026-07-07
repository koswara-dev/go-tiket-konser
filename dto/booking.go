package dto

import "github.com/google/uuid"

type BookingDetailRequest struct {
	TicketCategoryID uuid.UUID `json:"ticket_category_id" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,gte=1"`
}

type BookingRequest struct {
	CustomerID     uuid.UUID              `json:"customer_id" binding:"required"`
	BookingDetails []BookingDetailRequest `json:"booking_details" binding:"required"`
}
