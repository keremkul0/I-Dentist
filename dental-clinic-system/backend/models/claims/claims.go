package claims

import (
	"dental-clinic-system/models/user"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string       `json:"email"`
	Roles []*user.Role `json:"roles"` // Çoklu rol desteği
	jwt.RegisteredClaims
}
