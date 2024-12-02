package verifyEmail

import (
	"dental-clinic-system/models"
	"time"
)

type UserService interface {
	GetUserByEmail(email string) (models.UserGetModel, error)
}

type TokenService interface {
	AddTokenToBlacklistService(token string, expireTime time.Time)
	IsTokenBlacklistedService(token string) bool
}

type VerifyEmailHandler struct {
	userService  UserService
	tokenService TokenService
}
