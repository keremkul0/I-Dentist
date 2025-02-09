package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}
