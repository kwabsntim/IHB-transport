package auth

import "golang.org/x/crypto/bcrypt"

// creating the hashing functions for passwords
func HashPassword(password string) (string, error) {
	//implementations will be added here
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checking the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
