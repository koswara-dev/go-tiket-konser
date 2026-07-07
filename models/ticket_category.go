package models

import (
	"github.com/google/uuid"
)

type TicketCategory struct {
	BaseModel
	ConcertID      uuid.UUID `gorm:"type:uuid;not null" json:"concert_id" binding:"required"`
	Concert        *Concert  `gorm:"foreignKey:ConcertID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"concert,omitempty"`
	Name           string    `gorm:"type:varchar(50);not null" json:"name" binding:"required"`
	Price          float64   `gorm:"type:numeric(12,2);not null" json:"price" binding:"required,gt=0"`
	TotalQuota     int       `gorm:"not null" json:"total_quota" binding:"required,gte=0"`
	AvailableQuota int       `gorm:"not null" json:"available_quota" binding:"required,gte=0"`
}
