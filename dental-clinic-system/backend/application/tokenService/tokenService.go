package tokenService

import (
	"dental-clinic-system/repository/tokenRepository"
	"time"
)

type TokenService interface {
	DeleteExpiredTokensService()
	AddTokenToBlacklistService(token string, expireTime time.Time)
	IsTokenBlacklistedService(token string) bool
}

type tokenService struct {
	tokenRepository tokenRepository.TokenRepository
}

func NewTokenService(tokenRepository tokenRepository.TokenRepository) *tokenService {
	return &tokenService{
		tokenRepository: tokenRepository,
	}
}

func (s *tokenService) DeleteExpiredTokensService() {
	s.tokenRepository.DeleteExpiredTokensRepo()
}

func (s *tokenService) AddTokenToBlacklistService(token string, expireTime time.Time) {
	s.tokenRepository.AddTokenToBlacklistRepo(token, expireTime)
}

func (s *tokenService) IsTokenBlacklistedService(token string) bool {
	return s.tokenRepository.IsTokenBlacklistedRepo(token)
}
