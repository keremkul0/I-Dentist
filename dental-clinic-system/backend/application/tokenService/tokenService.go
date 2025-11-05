package tokenService

import (
	"context"
	"time"
)

type TokenRepository interface {
	DeleteExpiredTokens(ctx context.Context)
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

type tokenService struct {
	tokenRepository TokenRepository
}

func NewTokenService(tokenRepository TokenRepository) *tokenService {
	return &tokenService{
		tokenRepository: tokenRepository,
	}
}

func (s *tokenService) DeleteExpiredTokens(ctx context.Context) {
	s.tokenRepository.DeleteExpiredTokens(ctx)
}

func (s *tokenService) AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error {
	return s.tokenRepository.AddTokenToBlacklist(ctx, token, expireTime)
}
func (s *tokenService) IsTokenBlacklisted(ctx context.Context, token string) bool {
	return s.tokenRepository.IsTokenBlacklisted(ctx, token)
}
