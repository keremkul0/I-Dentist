package userService

import (
	"dental-clinic-system/helpers"
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
	users, err := s.userRepository.GetUsersRepo(ClinicID)
	if err != nil {
		return nil, err
	}

	var usersGetModel []models.UserGetModel

	for _, user := range users {
		usersGetModel = append(usersGetModel, helpers.UserConvertor(user))
	}

	return usersGetModel, nil
}

func (s *userService) GetUser(id uint) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUserRepo(id)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return helpers.UserConvertor(user), nil
}

func (s *userService) GetUserByEmail(email string) (models.UserGetModel, error) {
	user, err := s.userRepository.GetUserByEmailRepo(email)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return helpers.UserConvertor(user), nil
}

func (s *userService) CreateUser(user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.CreateUserRepo(user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return helpers.UserConvertor(user), nil
}

func (s *userService) UpdateUser(user models.User) (models.UserGetModel, error) {
	user, err := s.userRepository.UpdateUserRepo(user)
	if err != nil {
		return models.UserGetModel{}, err
	}

	return helpers.UserConvertor(user), nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUserRepo(id)
}

func (s *userService) CheckUserExist(user models.UserGetModel) bool {
	return s.userRepository.CheckUserExistRepo(user)
}
