package tokenService

import (
	"time"
)

type TokenRepository interface {
	DeleteExpiredTokensRepo()
	AddTokenToBlacklistRepo(token string, expireTime time.Time)
	IsTokenBlacklistedRepo(token string) bool
}

type tokenService struct {
	tokenRepository TokenRepository
}

func NewTokenService(tokenRepository TokenRepository) *tokenService {
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
