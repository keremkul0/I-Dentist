package loginService

import (
	"dental-clinic-system/models"
)

type LoginRepository interface {
	Login(email string, password string) (models.Login, error)
}

type loginService struct {
	loginRepository LoginRepository
}

func NewLoginService(loginRepository LoginRepository) *loginService {
	return &loginService{
		loginRepository: loginRepository,
	}
}

func (s *loginService) Login(email string, password string) (models.Login, error) {
	return s.loginRepository.Login(email, password)
}
