package dto

import "time"

type ConcertRequest struct {
	Title       string `json:"title" binding:"required,min=5,max=150"`
	Description string `json:"description" binding:"max=255"`
	Date        string `json:"date" binding:"required,future_date"`
	Venue       string `json:"venue" binding:"required,min=3,max=100"`
	Status      string `json:"status" binding:"required,oneof=active upcoming canceled"`
}

type ConcertResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	Venue       string    `json:"venue"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ConcertQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,numeric,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,numeric,gte=1"`
	Search string `form:"search"`
	Sort   string `form:"sort" binding:"omitempty,oneof=date_asc date_desc title_asc title_desc"`
}
