package signUpUserService

import (
	"context"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models/user"
	"dental-clinic-system/validations"
	"errors"
)

type UserRepository interface {
	CheckUserExist(ctx context.Context, userModel user.UserGetModel) (bool, error)
}

type RedisRepository interface {
	SetData(ctx context.Context, data any) (string, error)
}

type UserService interface {
	HashPassword(password string) string
}

type signUpUserService struct {
	userRepository  UserRepository
	redisRepository RedisRepository
	userService     UserService
}

func NewSignUpUserService(userRepository UserRepository, redisRepository RedisRepository, userService UserService) *signUpUserService {
	return &signUpUserService{
		userRepository:  userRepository,
		redisRepository: redisRepository,
		userService:     userService,
	}
}

func (s *signUpUserService) SignUpUser(ctx context.Context, user user.User) (string, error) {
	err := validations.UserValidation(&user)
	if err != nil {
		return "", errors.New("User validation errors")
	}

	user.Password = s.userService.HashPassword(user.Password)
	userGetModel := mapper.MapUserToUserGetModel(user)

	exists, err := s.userRepository.CheckUserExist(ctx, userGetModel)
	if err != nil {
		return "", errors.New("Error checking user existence")
	}
	if exists {
		return "", errors.New("User already exist")
	}

	userID, err := s.redisRepository.SetData(ctx, user)
	if err != nil {
		return "", errors.New("Cache user service errors")
	}

	return userID, nil
}
