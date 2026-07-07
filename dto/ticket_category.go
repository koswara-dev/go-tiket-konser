package dto

import (
	"time"

	"github.com/google/uuid"
)

type TicketCategoryRequest struct {
	ConcertID      uuid.UUID `json:"concert_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Price          float64   `json:"price" binding:"required,gt=400000"`
	TotalQuota     int       `json:"total_quota" binding:"required,gte=100"`
	AvailableQuota int       `json:"available_quota" binding:"required,gte=0"`
}

type TicketCategoryResponse struct {
	ID             uuid.UUID       `json:"id"`
	ConcertID      uuid.UUID       `json:"concert_id"`
	Concert        ConcertResponse `json:"concert"`
	Name           string          `json:"name"`
	Price          float64         `json:"price"`
	TotalQuota     int             `json:"total_quota"`
	AvailableQuota int             `json:"available_quota"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
