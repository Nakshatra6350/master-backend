package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if the given password matches the bcrypt hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
