package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	// Hashing the password (using bcrypt)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
