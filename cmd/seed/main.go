package main

//this creates a default admin account if not created
import (
	"ihb-transport/internal/auth"
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func seedAdmin(db *gorm.DB) {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		log.Println("Admin account exists, skipping seeding.")
		return
	}
	// Read from environment variable
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("no admin password found")
	}
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		log.Fatal("no admin email found")
	}
	//hashing password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
		return
	}
	//creating admin account
	admin := models.Admin{
		Email:        email,
		PasswordHash: hashedPassword,
	}
	db.Create(&admin)

}
func main() {
	godotenv.Load()
	database.Connection()
	db := database.DB

	seedAdmin(db)
}
