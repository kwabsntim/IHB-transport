package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateRequired checks if a required field is not empty
func ValidateRequired(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateMinLength checks if value meets minimum length
func ValidateMinLength(fieldName, value string, minLength int) error {
	if len(strings.TrimSpace(value)) < minLength {
		return errors.New(fieldName + " must be at least " + fmt.Sprintf("%d", minLength) + " characters")
	}
	return nil
}

// ValidateMaxLength checks if value doesn't exceed maximum length
func ValidateMaxLength(fieldName, value string, maxLength int) error {
	if len(value) > maxLength {
		return errors.New(fieldName + " cannot exceed " + fmt.Sprintf("%d", maxLength) + " characters")
	}
	return nil
}

// ValidatePositiveNumber checks if a number is positive
func ValidatePositiveNumber(fieldName string, value float64) error {
	if value < 0 {
		return errors.New(fieldName + " cannot be negative")
	}
	return nil
}

// IsValidEmail is a helper that returns bool instead of error
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}
