package loginService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/loginRepository"
)

type LoginService interface {
	Login(email string, password string) (models.Login, error)
}

type loginService struct {
	authRepository loginRepository.LoginRepository
}

func NewLoginService(authRepository loginRepository.LoginRepository) *loginService {
	return &loginService{
		authRepository: authRepository,
	}
}

func (s *loginService) Login(email string, password string) (models.Login, error) {
	return s.authRepository.Login(email, password)
}
