package helpers

import (
	"dental-clinic-system/models/claims"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

func GenerateJWTToken(email string, expirationTime time.Time) (string, error) {

	claims := &claims.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(GetJWTKey())
	if err != nil {
		log.Println("Could not create token")
		return "", err
	}
	return tokenString, nil
}
