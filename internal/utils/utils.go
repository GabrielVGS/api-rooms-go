package utils

import (
	"api-go/internal/server/dtos"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("Checking password hash:", err)
	if err != nil {
		return false
	}
	return true
}

func IsValidEmail(email string) bool {
	// A simple regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched := regexp.MustCompile(emailRegex).MatchString(email)
	return matched
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	errorResponse := dtos.ErrorResponse{
		Message: message,
		Status:  status,
	}

	json.NewEncoder(w).Encode(errorResponse)
}
