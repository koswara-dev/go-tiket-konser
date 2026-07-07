package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	CustomerID uuid.UUID `json:"customer_id,omitempty"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}
