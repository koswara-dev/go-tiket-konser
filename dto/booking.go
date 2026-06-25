package dto

type BookingDetailRequest struct {
	TicketCategoryID int `json:"ticket_category_id" binding:"required"`
	Quantity         int `json:"quantity" binding:"required,gte=1"`
}

type BookingRequest struct {
	CustomerID     int                    `json:"customer_id" binding:"required"`
	BookingDetails []BookingDetailRequest `json:"booking_details" binding:"required"`
}
