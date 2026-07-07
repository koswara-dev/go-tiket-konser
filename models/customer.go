package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name   string    `gorm:"not null" json:"name"`
	Email  string    `gorm:"not null;unique" json:"email"`
}
