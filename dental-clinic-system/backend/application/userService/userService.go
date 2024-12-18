package userService

import (
	"context"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models"
)

type UserRepository interface {
	GetUsers(ctx context.Context, ClinicID uint) ([]models.User, error)
	GetUser(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, user models.UserGetModel) bool
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers(ctx context.Context, ClinicID uint) ([]models.UserGetModel, error) {
	users, err := s.userRepository.GetUsers(ctx, ClinicID)
	if err != nil {
		return nil, err
	}

	var usersGetModel []models.UserGetModel

	for _, user := range users {
		usersGetModel = append(usersGetModel, mapper.UserMapper(user))
	}

	return usersGetModel, nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) CreateUser(ctx context.Context, user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepository.DeleteUser(ctx, id)
}

func (s *userService) CheckUserExist(ctx context.Context, user models.UserGetModel) bool {
	return s.userRepository.CheckUserExist(ctx, user)
}
