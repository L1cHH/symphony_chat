package utils

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Proves correctness of user input
func IsCorrectLoginFormat(login string) bool {
	if len(login) < 6 || len(login) > 30 {
		return false
	}

	disallowedSymbols := "#$%^&*!?~`':;_-+=<>,/|}{[]()"

	if strings.ContainsAny(login, disallowedSymbols) {
		return false
	} else {
		return true
	}

}

func IsCorrectPasswordFormat(password string) bool {
	if len(password) < 10 || len(password) > 25 {
		return false
	}

	for _, char := range password {
		// Проверяем, входит ли символ в диапазон ASCII (0–127)
		if char > unicode.MaxASCII {
			return false
		}
	}

	return true
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
