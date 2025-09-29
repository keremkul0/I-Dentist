package jwtService

import (
	"dental-clinic-system/models/claims"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	jwtSecret []byte
}

func NewJwtService(jwtSecret string) *jwtService {
	return &jwtService{
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *jwtService) GenerateJWTToken(email string, expirationTime time.Time) (string, error) {

	claims := &claims.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		print("Could not create token")
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GetJwtKey() []byte {
	return s.jwtSecret
}

func (s *jwtService) ParseToken(tokenStr string) (*claims.Claims, error) {
	// Token boş mu kontrol et
	if tokenStr == "" {
		return nil, errors.New("token is required")
	}

	// Claims yapısını oluştur
	claims := &claims.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.GetJwtKey(), nil // JWT anahtarını sağlayan fonksiyon
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

func (s *jwtService) ParseTokenFromCookie(r *http.Request) (*claims.Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return s.ParseToken(cookie.Value)
}
