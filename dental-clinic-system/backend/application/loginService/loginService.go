package loginService

import (
	"context"
	"dental-clinic-system/models/auth"
)

type LoginRepository interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
}

type loginService struct {
	loginRepository LoginRepository
}

func NewLoginService(loginRepository LoginRepository) *loginService {
	return &loginService{
		loginRepository: loginRepository,
	}
}

func (s *loginService) Login(ctx context.Context, email string, password string) (auth.Login, error) {
	return s.loginRepository.Login(ctx, email, password)
}
