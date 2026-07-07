package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlacklistedTokenRepository interface {
	BlacklistToken(token string, expiresAt time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
}

type blacklistedTokenRepository struct {
	redisClient *redis.Client
}

func NewBlacklistedTokenRepository(redisClient *redis.Client) BlacklistedTokenRepository {
	return &blacklistedTokenRepository{redisClient: redisClient}
}

func (r *blacklistedTokenRepository) BlacklistToken(token string, expiresAt time.Time) error {
	ctx := context.Background()
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	return r.redisClient.Set(ctx, "blacklist:"+token, "true", ttl).Err()
}

func (r *blacklistedTokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	val, err := r.redisClient.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
