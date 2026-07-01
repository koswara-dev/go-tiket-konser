package dto

import "time"

type CustomerUpdateRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CustomerResponse struct {
	ID        int       `json:"id"`
	UserID    uint      `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,numeric,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,numeric,gte=1"`
	Search string `form:"search"`
	Sort   string `form:"sort" binding:"omitempty,oneof=name_asc name_desc email_asc email_desc"`
}

