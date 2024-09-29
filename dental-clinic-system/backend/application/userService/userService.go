package userService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/userRepository"
)

type UserService interface {
	GetUsers(ClinicID uint) ([]models.UserGetModel, error)
	GetUser(id uint) (models.UserGetModel, error)
	GetUserByEmail(email string) (models.UserGetModel, error)
	CreateUser(user models.User) (models.UserGetModel, error)
	UpdateUser(user models.User) (models.UserGetModel, error)
	DeleteUser(id uint) error
	CheckUserExist(user models.UserGetModel) bool
}

type userService struct {
	userRepository userRepository.UserRepository
}

func NewUserService(userRepository userRepository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers(ClinicID uint) ([]models.UserGetModel, error) {
	return s.userRepository.GetUsersRepo(ClinicID)
}

func (s *userService) GetUser(id uint) (models.UserGetModel, error) {
	return s.userRepository.GetUserRepo(id)
}

func (s *userService) GetUserByEmail(email string) (models.UserGetModel, error) {
	return s.userRepository.GetUserByEmailRepo(email)
}

func (s *userService) CreateUser(user models.User) (models.UserGetModel, error) {
	return s.userRepository.CreateUserRepo(user)
}

func (s *userService) UpdateUser(user models.User) (models.UserGetModel, error) {
	return s.userRepository.UpdateUserRepo(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUserRepo(id)
}

func (s *userService) CheckUserExist(user models.UserGetModel) bool {
	return s.userRepository.CheckUserExistRepo(user)
}
