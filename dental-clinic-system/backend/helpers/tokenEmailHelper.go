package helpers

import (
	"dental-clinic-system/api/auth"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func TokenEmailHelper(r *http.Request) *auth.Claims {
	// Extract email from JWT
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}

	tokenStr := cookie.Value
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	return claims
}
