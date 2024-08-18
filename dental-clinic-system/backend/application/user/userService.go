package user

import (
	"dental-clinic-system/repository/user"
)

type UserService interface {
}

type userService struct {
	userRepository user.UserRepository
}

func NewAppointmentService(userRepository user.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}
