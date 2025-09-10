package userService

import (
	"context"
	"dental-clinic-system/mapper"
	modelerrors "dental-clinic-system/models/errors"

	"dental-clinic-system/models/user"
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrUserNotFound Error types
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	GetUsers(ctx context.Context, clinicID uint) ([]user.User, error)
	GetUser(ctx context.Context, id uint) (user.User, error)
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	CreateUser(ctx context.Context, usr user.User) (user.User, error)
	UpdateUser(ctx context.Context, usr user.User) (user.User, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, userModel user.UserGetModel) (bool, error)
}

// UserService handles user-related business logic
type UserService struct {
	userRepository UserRepository
	roleService    RoleService
}

type RoleService interface {
	UserHasRole(user user.UserGetModel, roleName string) bool
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo UserRepository, roleService RoleService) *UserService {
	return &UserService{
		userRepository: userRepo,
		roleService:    roleService,
	}
}

// GetUsers retrieves all users for a specific clinic and maps them to UserGetModel
func (s *UserService) GetUsers(ctx context.Context, clinicID uint) ([]user.UserGetModel, error) {
	log.Info().
		Str("operation", "GetUsers").
		Uint("clinic_id", clinicID).
		Msg("Fetching users for clinic")

	users, err := s.userRepository.GetUsers(ctx, clinicID)
	if err != nil {
		log.Error().
			Str("operation", "GetUsers").
			Err(err).
			Uint("clinic_id", clinicID).
			Msg("Failed to retrieve users")
		return nil, err
	}

	var usersGetModel []user.UserGetModel
	for _, usr := range users {
		usersGetModel = append(usersGetModel, mapper.MapUserToUserGetModel(usr))
	}

	log.Info().
		Str("operation", "GetUsers").
		Int("count", len(usersGetModel)).
		Msgf("Retrieved %d users successfully", len(usersGetModel))

	return usersGetModel, nil
}

// GetUser retrieves a single user by its ID and maps it to UserGetModel
func (s *UserService) GetUser(ctx context.Context, id uint) (user.UserGetModel, error) {
	log.Info().
		Str("operation", "GetUser").
		Uint("user_id", id).
		Msg("Fetching user by ID")

	usr, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetUser").
				Err(err).
				Uint("user_id", id).
				Msg("User not found")
			return user.UserGetModel{}, ErrUserNotFound
		}
		log.Error().
			Str("operation", "GetUser").
			Err(err).
			Uint("user_id", id).
			Msg("Failed to retrieve user")
		return user.UserGetModel{}, err
	}

	log.Info().
		Str("operation", "GetUser").
		Uint("user_id", id).
		Msg("User retrieved successfully")

	return mapper.MapUserToUserGetModel(usr), nil
}

// GetUserByEmail retrieves a single user by its email and maps it to UserGetModel
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error) {
	log.Info().
		Str("operation", "GetUserByEmail").
		Str("email", email).
		Msg("Fetching user by email")

	usr, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetUserByEmail").
				Err(err).
				Str("email", email).
				Msg("User not found")
			return user.UserGetModel{}, ErrUserNotFound
		}
		log.Error().
			Str("operation", "GetUserByEmail").
			Err(err).
			Str("email", email).
			Msg("Failed to retrieve user by email")
		return user.UserGetModel{}, err
	}

	log.Info().
		Str("operation", "GetUserByEmail").
		Str("email", email).
		Msg("User retrieved by email successfully")

	return mapper.MapUserToUserGetModel(usr), nil
}

// CreateUser creates a new user record and maps it to UserGetModel
func (s *UserService) CreateUser(ctx context.Context, newUsr user.User) (user.UserGetModel, error) {
	log.Info().
		Str("operation", "CreateUser").
		Str("email", newUsr.Email).
		Msg("Creating new user")

	usr, err := s.userRepository.CreateUser(ctx, newUsr)
	if err != nil {
		log.Error().
			Str("operation", "CreateUser").
			Err(err).
			Msg("Failed to create user")
		return user.UserGetModel{}, err
	}

	log.Info().
		Str("operation", "CreateUser").
		Uint("user_id", usr.ID).
		Msg("User created successfully")

	return mapper.MapUserToUserGetModel(usr), nil
}

func (s *UserService) CreateUserWithAuthorization(ctx context.Context, newUser user.User, authUserEmail string) (user.UserGetModel, error) {
	authenticatedUser, err := s.GetUserByEmail(ctx, authUserEmail)
	if err != nil {
		return user.UserGetModel{}, err
	}

	if !s.canCreateUser(authenticatedUser, newUser) {
		return user.UserGetModel{}, &modelerrors.UnauthorizedError{Message: "Insufficient permissions"}
	}

	tempUserGetModel := mapper.MapUserToUserGetModel(newUser)
	if exists, _ := s.CheckUserExist(ctx, tempUserGetModel); exists {
		return user.UserGetModel{}, &modelerrors.ValidationError{Message: "User already exists"}
	}

	newUser.Password = s.HashPassword(newUser.Password)

	return s.CreateUser(ctx, newUser)
}

func (s *UserService) canCreateUser(authenticatedUser user.UserGetModel, newUser user.User) bool {
	if authenticatedUser.ClinicID != newUser.ClinicID {
		return false
	}
	return s.roleService.UserHasRole(authenticatedUser, "Clinic Admin") ||
		s.roleService.UserHasRole(authenticatedUser, "Superadmin")
}

func (s *UserService) HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// UpdateUser updates an existing user record and maps it to UserGetModel
func (s *UserService) UpdateUser(ctx context.Context, updatedUsr user.User) (user.UserGetModel, error) {
	log.Info().
		Str("operation", "UpdateUser").
		Uint("user_id", updatedUsr.ID).
		Msg("Updating user")

	usr, err := s.userRepository.UpdateUser(ctx, updatedUsr)
	if err != nil {
		log.Error().
			Str("operation", "UpdateUser").
			Err(err).
			Uint("user_id", updatedUsr.ID).
			Msg("Failed to update user")
		return user.UserGetModel{}, err
	}

	log.Info().
		Str("operation", "UpdateUser").
		Uint("user_id", usr.ID).
		Msg("User updated successfully")

	return mapper.MapUserToUserGetModel(usr), nil
}

// DeleteUser deletes a user record by its ID
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	log.Info().
		Str("operation", "DeleteUser").
		Uint("user_id", id).
		Msg("Deleting user")

	err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		log.Error().
			Str("operation", "DeleteUser").
			Err(err).
			Uint("user_id", id).
			Msg("Failed to delete user")
		return err
	}

	log.Info().
		Str("operation", "DeleteUser").
		Uint("user_id", id).
		Msg("User deleted successfully")

	return nil
}

// CheckUserExist checks if a user exists based on national ID, email, or phone number
func (s *UserService) CheckUserExist(ctx context.Context, userModel user.UserGetModel) (bool, error) {
	log.Info().
		Str("operation", "CheckUserExist").
		Str("email", userModel.Email).
		Str("national_id", userModel.NationalID).
		Str("phone_number", userModel.PhoneNumber).
		Msg("Checking if user exists")

	exists, err := s.userRepository.CheckUserExist(ctx, userModel)
	if err != nil {
		log.Error().
			Str("operation", "CheckUserExist").
			Err(err).
			Msg("Failed to check user existence")
		return false, err
	}

	if exists {
		log.Info().
			Str("operation", "CheckUserExist").
			Msg("User exists")
	} else {
		log.Info().
			Str("operation", "CheckUserExist").
			Msg("User does not exist")
	}

	return exists, nil
}
