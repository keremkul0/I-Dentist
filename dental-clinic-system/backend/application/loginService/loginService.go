package loginService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/loginRepository"
)

type LoginService interface {
	Login(email string, password string) (models.Login, error)
}

type loginService struct {
	loginRepository loginRepository.LoginRepository
}

func NewLoginService(loginRepository loginRepository.LoginRepository) *loginService {
	return &loginService{
		loginRepository: loginRepository,
	}
}

func (s *loginService) Login(email string, password string) (models.Login, error) {
	return s.loginRepository.Login(email, password)
}
