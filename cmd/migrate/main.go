package main

import (
	"IHB-transport/internal/database"
	"IHB-transport/migrations"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Connecting to database for migration...")
	// Connect to Postgres
	database.Connection()
	log.Println("Database connection successful.")

	log.Println("Starting database migration...")
	// Auto-migrate models
	migrations.Migrate()
	log.Println("Migration process completed.")
}
