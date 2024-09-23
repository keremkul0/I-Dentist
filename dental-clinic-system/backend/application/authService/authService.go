package authService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/authRepository"
)

type AuthService interface {
	Login(email string, password string) (models.Auth, error)
}

type authService struct {
	authRepository authRepository.AuthRepository
}

func NewAuthService(authRepository authRepository.AuthRepository) *authService {
	return &authService{
		authRepository: authRepository,
	}
}

func (s *authService) Login(email string, password string) (models.Auth, error) {
	return s.authRepository.Login(email, password)
}
