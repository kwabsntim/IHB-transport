package main

import (
	"ihb-transport/internal/database"
	"ihb-transport/migrations"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file FIRST
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
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
