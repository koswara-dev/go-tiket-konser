package models

import "time"

type Concert struct {
	BaseModel
	Title        string    `gorm:"not null" json:"title" binding:"required"`
	Description  string    `gorm:"type:text" json:"description"`
	Date         time.Time `gorm:"not null" json:"date" binding:"required"`
	Venue        string    `gorm:"not null" json:"venue" binding:"required"`
	Status       string    `gorm:"default:active" json:"status"`
	PosterURL    string    `gorm:"null" json:"poster_url" binding:"required"`
	ThumbnailURL string    `gorm:"null" json:"thumbnail_url" binding:"required"`
	RulesPDFURL  string    `gorm:"null" json:"rules_pdf_url" binding:"required"`
}
