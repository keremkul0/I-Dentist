package userService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/userRepository"
)

type UserService interface {
	GetUsers() ([]models.User, error)
	GetUser(id uint) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepository userRepository.UserRepository
}

func NewUserService(userRepository userRepository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers() ([]models.User, error) {
	return s.userRepository.GetUsers()
}

func (s *userService) GetUser(id uint) (models.User, error) {
	return s.userRepository.GetUser(id)
}

func (s *userService) CreateUser(user models.User) (models.User, error) {
	return s.userRepository.CreateUser(user)
}

func (s *userService) UpdateUser(user models.User) (models.User, error) {
	return s.userRepository.UpdateUser(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUser(id)
}
