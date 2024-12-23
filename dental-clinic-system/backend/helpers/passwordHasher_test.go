package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword := HashPassword(password)

	// Check if the hashed password is not equal to the original password
	if hashedPassword == password {
		t.Errorf("Hashed password should not be equal to the original password")
	}

	// Check if the hashed password can be verified using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("Hashed password could not be verified: %v", err)
	}
}
