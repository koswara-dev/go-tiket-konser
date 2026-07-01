package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName     string    `gorm:"not null" json:"full_name"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Password     string    `gorm:"not null" json:"password"`
	Role         string    `gorm:"not null;default:'customer'" json:"role"`
	Customer     *Customer `gorm:"foreignKey:UserID" json:"customer,omitempty"`
	OtpCode      string    `json:"otp_code,omitempty"`
	OtpExpiredAt time.Time `json:"otp_expired_at,omitempty"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
}
