package userRepository

import (
	"context"
	"dental-clinic-system/models/user"
	"errors"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles user-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetUsers retrieves all users for a specific clinic
func (repo *Repository) GetUsers(ctx context.Context, clinicID uint) ([]user.User, error) {
	var usersList []user.User
	result := repo.DB.WithContext(ctx).Where("clinic_id = ?", clinicID).Find(&usersList)
	if result.Error != nil {
		log.Error().
			Str("operation", "GetUsers").
			Err(result.Error).
			Uint("clinic_id", clinicID).
			Msg("Failed to retrieve users")
		return nil, result.Error
	}
	log.Info().
		Str("operation", "GetUsers").
		Uint("clinic_id", clinicID).
		Int("count", len(usersList)).
		Msg("Retrieved users successfully")
	return usersList, nil
}

// GetUser retrieves a single user by its ID
func (repo *Repository) GetUser(ctx context.Context, id uint) (user.User, error) {
	var usr user.User
	result := repo.DB.WithContext(ctx).Where("id = ?", id).First(&usr)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().
				Str("operation", "GetUser").
				Err(result.Error).
				Uint("user_id", id).
				Msg("User not found")
		} else {
			log.Error().
				Str("operation", "GetUser").
				Err(result.Error).
				Uint("user_id", id).
				Msg("Failed to retrieve user")
		}
		return user.User{}, result.Error
	}
	log.Info().
		Str("operation", "GetUser").
		Uint("user_id", id).
		Msg("Retrieved user successfully")
	return usr, nil
}

// GetUserByEmail retrieves a single user by its email
func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	var usr user.User
	result := repo.DB.WithContext(ctx).Where("email = ?", email).First(&usr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetUserByEmail").
				Err(result.Error).
				Str("email", email).
				Msg("User not found")
		} else {
			log.Error().
				Str("operation", "GetUserByEmail").
				Err(result.Error).
				Str("email", email).
				Msg("Failed to retrieve user by email")
		}
		return user.User{}, result.Error
	}
	log.Info().
		Str("operation", "GetUserByEmail").
		Str("email", email).
		Msg("Retrieved user by email successfully")
	return usr, nil
}

// CreateUser creates a new user record in the database
func (repo *Repository) CreateUser(ctx context.Context, newUser user.User) (user.User, error) {
	result := repo.DB.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreateUser").
			Err(result.Error).
			Msg("Failed to create user")
		return user.User{}, result.Error
	}
	log.Info().
		Str("operation", "CreateUser").
		Uint("user_id", newUser.ID).
		Msg("User created successfully")
	return newUser, nil
}

// UpdateUser updates an existing user record in the database
func (repo *Repository) UpdateUser(ctx context.Context, updatedUser user.User) (user.User, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedUser)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdateUser").
			Err(result.Error).
			Uint("user_id", updatedUser.ID).
			Msg("Failed to update user")
		return user.User{}, result.Error
	}
	log.Info().
		Str("operation", "UpdateUser").
		Uint("user_id", updatedUser.ID).
		Msg("User updated successfully")
	return updatedUser, nil
}

// DeleteUser deletes a user record from the database by its ID
func (repo *Repository) DeleteUser(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&user.User{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteUser").
			Err(result.Error).
			Uint("user_id", id).
			Msg("Failed to delete user")
		return result.Error
	}
	log.Info().
		Str("operation", "DeleteUser").
		Uint("user_id", id).
		Msg("User deleted successfully")
	return nil
}

// CheckUserExist checks if a user exists based on national ID, email, or phone number
func (repo *Repository) CheckUserExist(ctx context.Context, userModel user.UserGetModel) (bool, error) {
	var count int64
	result := repo.DB.WithContext(ctx).
		Model(&user.User{}).
		Where("national_id = ? OR email = ? OR phone_number = ?", userModel.NationalID, userModel.Email, userModel.PhoneNumber).
		Count(&count)

	if result.Error != nil {
		log.Error().
			Str("operation", "CheckUserExist").
			Err(result.Error).
			Msg("Failed to check user existence")
		return false, result.Error
	}

	exists := count > 0
	if exists {
		log.Info().
			Str("operation", "CheckUserExist").
			Msgf("User exists with National ID: %s, Email: %s, Phone Number: %s", userModel.NationalID, userModel.Email, userModel.PhoneNumber)
	} else {
		log.Info().
			Str("operation", "CheckUserExist").
			Msg("User does not exist")
	}

	return exists, nil
}
