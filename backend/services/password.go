package services

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func VerifyPassword(password, storedValue string) (matched bool, needsRehash bool) {
	if CheckPasswordHash(password, storedValue) {
		return true, false
	}

	if storedValue == password {
		return true, true
	}

	return false, false
}

func IsBcryptHash(value string) bool {
	return strings.HasPrefix(value, "$2a$") || strings.HasPrefix(value, "$2b$") || strings.HasPrefix(value, "$2y$")
}
