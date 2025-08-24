package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default cost for bcrypt hashing
	DefaultCost = 12
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a plain text password with a hashed password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsValidPassword checks if a password meets the minimum requirements
func IsValidPassword(password string) bool {
	// Minimum 6 characters
	if len(password) < 6 {
		return false
	}
	
	// Add more validation rules as needed
	// For example: require uppercase, lowercase, numbers, special characters
	
	return true
}