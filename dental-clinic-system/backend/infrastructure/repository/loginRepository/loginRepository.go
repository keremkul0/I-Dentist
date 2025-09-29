package loginRepository

import (
	"context"
	"dental-clinic-system/models/auth"
	"dental-clinic-system/models/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles login-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// Login authenticates a user by email and password
func (repo *Repository) Login(ctx context.Context, email string, password string) (auth.Login, error) {
	var usr user.User
	result := repo.DB.WithContext(ctx).Where("email = ?", email).First(&usr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().Str("email", email).Msg("Login attempt with non-existent email")
			return auth.Login{}, result.Error
		}
		log.Error().Err(result.Error).Str("email", email).Msg("Failed to retrieve user during login")
		return auth.Login{}, result.Error
	}

	// Compare the hashed password with the plain password
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Warn().Str("email", email).Msg("Incorrect password attempt")
			return auth.Login{}, err
		}
		log.Error().Err(err).Str("email", email).Msg("Error comparing passwords")
		return auth.Login{}, err
	}

	log.Info().Str("email", email).Msg("User authenticated successfully")

	// Return only necessary information, avoiding sensitive data like Password
	return auth.Login{
		Email: usr.Email,
		// Diğer gerekli alanlar buraya eklenebilir (örn. Token)
	}, nil
}
