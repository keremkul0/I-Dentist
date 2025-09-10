package tokenRepository

import (
	"context"
	"errors"
	"time"

	"dental-clinic-system/models/token"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles token-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// DeleteExpiredTokens removes all expired tokens from the database
func (repo *Repository) DeleteExpiredTokens(ctx context.Context) {
	result := repo.DB.WithContext(ctx).
		Unscoped().
		Where("expire_time < ?", time.Now().Unix()).
		Delete(&token.ExpiredTokens{})

	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteExpiredTokens").
			Err(result.Error).
			Msg("Failed to delete expired tokens")
		return
	}

	log.Info().
		Str("operation", "DeleteExpiredTokens").
		Int64("deleted_count", result.RowsAffected).
		Msg("Expired tokens deleted successfully")
}

// AddTokenToBlacklist adds a token to the blacklist with an expiration time
func (repo *Repository) AddTokenToBlacklist(ctx context.Context, tokenStr string, expireTime time.Time) error {
	newToken := token.ExpiredTokens{
		Token:      tokenStr,
		ExpireTime: expireTime.Unix(),
	}

	result := repo.DB.WithContext(ctx).Create(&newToken)
	if result.Error != nil {
		log.Error().
			Str("operation", "AddTokenToBlacklist").
			Err(result.Error).
			Msg("Failed to add token to blacklist")
		return errors.New("redis set errors")
	}

	log.Info().
		Str("operation", "AddTokenToBlacklist").
		Str("token", tokenStr).
		Int64("expire_time", newToken.ExpireTime).
		Msg("Token added to blacklist successfully")
	return nil
}

// IsTokenBlacklisted checks if a token is present in the blacklist
func (repo *Repository) IsTokenBlacklisted(ctx context.Context, tokenStr string) bool {
	var count int64
	result := repo.DB.WithContext(ctx).
		Model(&token.ExpiredTokens{}).
		Where("token = ?", tokenStr).
		Count(&count)

	if result.Error != nil {
		log.Error().
			Str("operation", "IsTokenBlacklisted").
			Err(result.Error).
			Str("token", tokenStr).
			Msg("Failed to check if token is blacklisted")
		return false
	}

	if count > 0 {
		log.Info().
			Str("operation", "IsTokenBlacklisted").
			Str("token", tokenStr).
			Msg("Token is blacklisted")
	} else {
		log.Info().
			Str("operation", "IsTokenBlacklisted").
			Str("token", tokenStr).
			Msg("Token is not blacklisted")
	}

	return count > 0
}
