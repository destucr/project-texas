package utils

import (
    "golang.org/x/crypto/bcrypt"
    "fmt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain-text password
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
