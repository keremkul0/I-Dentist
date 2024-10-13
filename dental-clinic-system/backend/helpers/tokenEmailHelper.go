package helpers

import (
	"dental-clinic-system/models"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func TokenEmailHelper(r *http.Request) *models.Claims {
	// Extract email from JWT
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}

	tokenStr := cookie.Value
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	return claims
}
