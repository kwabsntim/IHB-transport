package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// function to get secret key from environment variable
func getSecretKey() []byte {
	secret_key := os.Getenv("JWT_SECRET_KEY")
	if secret_key == "" {
		log.Fatal("JWT_SECRET_KEY not found in environment")
	}
	return []byte(secret_key)
}

// creating JWT token
func CreateToken(email, role string) (string, error) {
	secretKey := getSecretKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"role":  role,
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	secretKey := getSecretKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	}
	return nil, fmt.Errorf("Inavlid token claims")
}
