package helpers

import (
	"dental-clinic-system/models/claims"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func TokenEmailHelper(tokenStr string) (*claims.Claims, error) {
	// Token boş mu kontrol et
	if tokenStr == "" {
		return nil, errors.New("token is required")
	}

	// Claims yapısını oluştur
	claims := &claims.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil // JWT anahtarını sağlayan fonksiyon
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
