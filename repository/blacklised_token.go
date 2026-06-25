package repository

import (
	"go-tiket-konser/models"
	"time"

	"gorm.io/gorm"
)

type BlacklistedTokenRepository interface {
	BlacklistToken(token string, expiresAt time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
}

type blacklistedTokenRepository struct {
	db *gorm.DB
}

func NewBlacklistedTokenRepository(db *gorm.DB) BlacklistedTokenRepository {
	return &blacklistedTokenRepository{db: db}
}

func (r *blacklistedTokenRepository) BlacklistToken(token string, expiresAt time.Time) error {
	return r.db.Create(&models.BlacklistedToken{
		TokenString: token,
		ExpiredAt:   expiresAt,
	}).Error
}

func (r *blacklistedTokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlacklistedToken{}).
		Where("token_string = ? AND expired_at > ?", token, time.Now()).
		Count(&count).Error
	return count > 0, err
}
