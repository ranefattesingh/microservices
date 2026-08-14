package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword encrypts a password using bcrypt hashing algorithm
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a hashed password with a plain text password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
