package authMiddleware

import (
	"dental-clinic-system/application/tokenService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/tokenRepository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
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

		tokenStr := cookie.Value
		dsn := os.Getenv("DNS")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newTokenRepository := tokenRepository.NewTokenRepository(db)
		tokenServiceClient := tokenService.NewTokenService(newTokenRepository)
		if tokenServiceClient.IsTokenBlacklistedService(tokenStr) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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
