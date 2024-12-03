package helpers

import (
	"dental-clinic-system/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func CookieTokenEmailHelper(r *http.Request) (*models.Claims, error) {
	// Extract email from JWT
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	tokenStr := cookie.Value
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
