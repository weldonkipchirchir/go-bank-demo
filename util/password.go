package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashedPassword returns the bcrypt hash of the password
func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash the password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
