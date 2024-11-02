package singUpUserService

import (
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/userRepository"
	"dental-clinic-system/validations"
	"errors"
)

type SignUpUserService interface {
	SignUpUserService(user models.User) (models.User, error)
}

type signUpUserService struct {
	userRepository userRepository.UserRepository
}

func NewSignUpUserService(userRepository userRepository.UserRepository) *signUpUserService {
	return &signUpUserService{
		userRepository: userRepository,
	}
}

func (s *signUpUserService) SignUpUserService(user models.User) (models.User, error) {

	err := validations.UserValidation(&user)

	if err != nil {
		return models.User{}, err
	}

	user.Password = helpers.HashPassword(user.Password)
	userGetModel := helpers.UserConvertor(user)

	if s.userRepository.CheckUserExistRepo(userGetModel) {
		return models.User{}, errors.New("User already exists")
	}

	return user, nil
}
