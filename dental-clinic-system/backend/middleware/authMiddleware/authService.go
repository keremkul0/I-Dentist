package authMiddleware

import (
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type TokenService interface {
	DeleteExpiredTokensService()
	AddTokenToBlacklistService(token string, expireTime time.Time)
	IsTokenBlacklistedService(token string) bool
}
type AuthMiddleware struct {
	TokenService TokenService
}

func NewAuthMiddleware(tokenService TokenService) *AuthMiddleware {
	return &AuthMiddleware{TokenService: tokenService}
}

func (auth *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if errors.Is(http.ErrNoCookie, err) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if auth.TokenService.IsTokenBlacklistedService(cookie.Value) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return helpers.GetJWTKey(), nil
		})

		if err != nil {
			if errors.Is(jwt.ErrSignatureInvalid, err) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
