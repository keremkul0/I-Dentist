package signUpClinicService

import (
	"context"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"
	"dental-clinic-system/validations"
	"errors"

	"github.com/rs/zerolog/log"
)

// Error types
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExist   = errors.New("user already exists")
	ErrClinicAlreadyExist = errors.New("clinic already exists")
	ErrClinicValidation   = errors.New("clinic validation error")
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, usr user.User) (user.User, error)
	CheckUserExist(ctx context.Context, userModel user.UserGetModel) (bool, error)
}

// ClinicRepository defines the interface for clinic-related database operations
type ClinicRepository interface {
	CreateClinic(ctx context.Context, cln clinic.Clinic) (clinic.Clinic, error)
	CheckClinicExist(ctx context.Context, cln clinic.Clinic) (bool, error)
}

// RedisRepository defines the interface for Redis-related operations
type RedisRepository interface {
	GetData(ctx context.Context, cacheKey string, target any) error
	DeleteData(ctx context.Context, cacheKey string) error
}

// SignUpClinicService handles the business logic for signing up clinics
type SignUpClinicService struct {
	clinicRepository ClinicRepository
	userRepository   UserRepository
	redisRepository  RedisRepository
}

// NewSignUpClinicService creates a new instance of SignUpClinicService
func NewSignUpClinicService(clinicRepo ClinicRepository, userRepo UserRepository, redisRepo RedisRepository) *SignUpClinicService {
	return &SignUpClinicService{
		clinicRepository: clinicRepo,
		userRepository:   userRepo,
		redisRepository:  redisRepo,
	}
}

// SignUpClinic registers a new clinic and its associated user
func (s *SignUpClinicService) SignUpClinic(ctx context.Context, cln clinic.Clinic, userCacheKey string) (clinic.Clinic, user.UserGetModel, error) {
	log.Info().
		Str("operation", "SignUpClinic").
		Str("user_cache_key", userCacheKey).
		Msg("Starting clinic signup process")

	// Retrieve user data from Redis cache
	var usr user.User
	err := s.redisRepository.GetData(ctx, userCacheKey, &usr)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Str("user_cache_key", userCacheKey).
			Msg("Failed to retrieve user from cache")
		return clinic.Clinic{}, user.UserGetModel{}, ErrUserNotFound
	}

	// Map user to UserGetModel for existence check
	userGetModel := mapper.MapUserToUserGetModel(usr)

	// Check if user already exists
	exists, err := s.userRepository.CheckUserExist(ctx, userGetModel)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Msg("Error while checking if user exists")
		return clinic.Clinic{}, user.UserGetModel{}, err
	}
	if exists {
		log.Warn().
			Str("operation", "SignUpClinic").
			Msg("User already exists")
		return clinic.Clinic{}, user.UserGetModel{}, ErrUserAlreadyExist
	}

	// Check if clinic already exists
	clinicExists, err := s.clinicRepository.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Msg("Error while checking if clinic exists")
		return clinic.Clinic{}, user.UserGetModel{}, err
	}
	if clinicExists {
		log.Warn().
			Str("operation", "SignUpClinic").
			Msg("Clinic already exists")
		return clinic.Clinic{}, user.UserGetModel{}, ErrClinicAlreadyExist
	}

	// Validate clinic data
	err = validations.ClinicValidation(&cln)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Msg("Clinic validation failed")
		return clinic.Clinic{}, user.UserGetModel{}, ErrClinicValidation
	}

	// Create clinic record in the database
	createdCln, err := s.clinicRepository.CreateClinic(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Msg("Failed to create clinic")
		return clinic.Clinic{}, user.UserGetModel{}, err
	}
	log.Info().
		Str("operation", "SignUpClinic").
		Uint("clinic_id", createdCln.ID).
		Msg("Clinic created successfully")

	// Associate user with the created clinic
	usr.ClinicID = createdCln.ID
	userGetModel.ClinicID = usr.ClinicID

	// Create user record in the database
	createdUsr, err := s.userRepository.CreateUser(ctx, usr)
	if err != nil {
		log.Error().
			Str("operation", "SignUpClinic").
			Err(err).
			Uint("clinic_id", createdCln.ID).
			Msg("Failed to create user")
		return createdCln, user.UserGetModel{}, err
	}
	log.Info().
		Str("operation", "SignUpClinic").
		Uint("user_id", createdUsr.ID).
		Msg("User created successfully")

	// Map the created user to UserGetModel
	createdUserGetModel := mapper.MapUserToUserGetModel(createdUsr)

	// Delete user data from Redis cache
	err = s.redisRepository.DeleteData(ctx, userCacheKey)
	if err != nil {
		log.Warn().
			Str("operation", "SignUpClinic").
			Err(err).
			Str("user_cache_key", userCacheKey).
			Msg("Failed to delete user cache")
	}

	return createdCln, createdUserGetModel, nil
}
