package main

import (
	
	"ihb-transport/internal/database"
	"ihb-transport/migrations"
	"log"
)

func main() {
	log.Println("Connecting to database for migration...")
	// Connect to Postgres
	database.Connection()
	log.Println("Database connection successful.")

	log.Println("Starting database migration...")
	// Auto-migrate models
	migrations.Migrate()
	log.Println("Migration process completed.")
}
