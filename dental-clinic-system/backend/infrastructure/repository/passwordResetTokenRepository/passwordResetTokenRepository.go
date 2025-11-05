package passwordResetTokenRepository

import (
	"context"
	"crypto/rand"
	"dental-clinic-system/models/token"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles password reset token-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// CreatePasswordResetToken creates a new password reset token for the given email
func (repo *Repository) CreatePasswordResetToken(ctx context.Context, email string) (token.PasswordResetToken, error) {
	// Generate random token
	tokenStr, err := generateRandomToken()
	if err != nil {
		log.Error().
			Str("operation", "CreatePasswordResetToken").
			Err(err).
			Str("email", email).
			Msg("Failed to generate random token")

		return token.PasswordResetToken{}, err
	}

	// Create token record
	resetToken := token.PasswordResetToken{
		Token:     tokenStr,
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 1), // 1 hour expiration
		IsUsed:    false,
	}

	result := repo.DB.WithContext(ctx).Create(&resetToken)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreatePasswordResetToken").
			Err(result.Error).
			Str("email", email).
			Msg("Failed to create password reset token")
		return token.PasswordResetToken{}, result.Error
	}

	log.Info().
		Str("operation", "CreatePasswordResetToken").
		Str("email", email).
		Uint("token_id", resetToken.ID).
		Msg("Password reset token created successfully")

	return resetToken, nil
}

// ValidatePasswordResetToken validates a password reset token
func (repo *Repository) ValidatePasswordResetToken(ctx context.Context, tokenStr string, email string) (token.PasswordResetToken, error) {
	var resetToken token.PasswordResetToken

	result := repo.DB.WithContext(ctx).Where(
		"token = ? AND email = ? AND is_used = false AND expires_at > ?",
		tokenStr, email, time.Now(),
	).First(&resetToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "ValidatePasswordResetToken").
				Str("email", email).
				Msg("Password reset token not found or expired")
		} else {
			log.Error().
				Str("operation", "ValidatePasswordResetToken").
				Err(result.Error).
				Str("email", email).
				Msg("Failed to validate password reset token")
		}
		return token.PasswordResetToken{}, result.Error
	}

	log.Info().
		Str("operation", "ValidatePasswordResetToken").
		Str("email", email).
		Uint("token_id", resetToken.ID).
		Msg("Password reset token validated successfully")

	return resetToken, nil
}

// MarkTokenAsUsed marks a password reset token as used
func (repo *Repository) MarkTokenAsUsed(ctx context.Context, tokenID uint) error {
	result := repo.DB.WithContext(ctx).Model(&token.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("is_used", true)

	if result.Error != nil {
		log.Error().
			Str("operation", "MarkTokenAsUsed").
			Err(result.Error).
			Uint("token_id", tokenID).
			Msg("Failed to mark token as used")
		return result.Error
	}

	log.Info().
		Str("operation", "MarkTokenAsUsed").
		Uint("token_id", tokenID).
		Msg("Token marked as used successfully")

	return nil
}

// DeleteExpiredTokens deletes all expired tokens from the database
func (repo *Repository) DeleteExpiredTokens(ctx context.Context) error {
	result := repo.DB.WithContext(ctx).Where("expires_at < ?", time.Now()).
		Delete(&token.PasswordResetToken{})

	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteExpiredTokens").
			Err(result.Error).
			Msg("Failed to delete expired tokens")
		return result.Error
	}

	log.Info().
		Str("operation", "DeleteExpiredTokens").
		Int64("deleted_count", result.RowsAffected).
		Msg("Expired tokens deleted successfully")

	return nil
}

// GetTokensByEmail retrieves all active tokens for a specific email
func (repo *Repository) GetTokensByEmail(ctx context.Context, email string) ([]token.PasswordResetToken, error) {
	var tokens []token.PasswordResetToken

	result := repo.DB.WithContext(ctx).Where(
		"email = ? AND is_used = false AND expires_at > ?",
		email, time.Now(),
	).Find(&tokens)

	if result.Error != nil {
		log.Error().
			Str("operation", "GetTokensByEmail").
			Err(result.Error).
			Str("email", email).
			Msg("Failed to retrieve tokens by email")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "GetTokensByEmail").
		Str("email", email).
		Int("count", len(tokens)).
		Msg("Retrieved tokens by email successfully")

	return tokens, nil
}

// InvalidateAllTokensForEmail marks all tokens for a specific email as used
func (repo *Repository) InvalidateAllTokensForEmail(ctx context.Context, email string) error {
	result := repo.DB.WithContext(ctx).Model(&token.PasswordResetToken{}).
		Where("email = ? AND is_used = false", email).
		Update("is_used", true)

	if result.Error != nil {
		log.Error().
			Str("operation", "InvalidateAllTokensForEmail").
			Err(result.Error).
			Str("email", email).
			Msg("Failed to invalidate tokens for email")
		return result.Error
	}

	log.Info().
		Str("operation", "InvalidateAllTokensForEmail").
		Str("email", email).
		Int64("updated_count", result.RowsAffected).
		Msg("All tokens invalidated for email successfully")

	return nil
}

// generateRandomToken generates a secure random token
func generateRandomToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
