package utils

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Proves correctness of user input
func IsCorrectFormat(input string) bool {
	if len(input) < 8 || len(input) > 16 {
		return false
	}

	disallowedSymbols := "#$%^&*!?~`':;_-+=<>,/|}{[]()"

	if !strings.ContainsAny(input, disallowedSymbols) {
		return true
	} else {
		return false
	}

}

// Hashing password
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(inputPassword string, correctPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(inputPassword))
	return err == nil
}
