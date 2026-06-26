package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"not null;unique" json:"email"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
