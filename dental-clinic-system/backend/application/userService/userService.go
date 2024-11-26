package userService

import (
	"dental-clinic-system/mapper"
	"dental-clinic-system/models"
)

type UserRepository interface {
	GetUsersRepo(ClinicID uint) ([]models.User, error)
	GetUserRepo(id uint) (models.User, error)
	GetUserByEmailRepo(email string) (models.User, error)
	CreateUserRepo(user models.User) (models.User, error)
	UpdateUserRepo(user models.User) (models.User, error)
	DeleteUserRepo(id uint) error
	CheckUserExistRepo(user models.UserGetModel) bool
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers(ClinicID uint) ([]models.UserGetModel, error) {
	users, err := s.userRepository.GetUsersRepo(ClinicID)
	if err != nil {
		return nil, err
	}

	var usersGetModel []models.UserGetModel

	for _, user := range users {
		usersGetModel = append(usersGetModel, mapper.UserMapper(user))
	}

	return usersGetModel, nil
}

func (s *userService) GetUser(id uint) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUserRepo(id)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) GetUserByEmail(email string) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUserByEmailRepo(email)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) CreateUser(user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.CreateUserRepo(user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) UpdateUser(user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.UpdateUserRepo(user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return mapper.UserMapper(user), nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUserRepo(id)
}

func (s *userService) CheckUserExist(user models.UserGetModel) bool {
	return s.userRepository.CheckUserExistRepo(user)
}
