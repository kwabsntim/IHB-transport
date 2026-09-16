package database

//this  has the database connections

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Missing required environment variable: %s", key)
	}
	return value
}

func Connection() {
	// Check if DATABASE_URL is provided (Supabase/Railway/Render connection string)
	databaseURL := os.Getenv("DATABASE_URL")

	var dsn string
	if databaseURL != "" {
		// Use full connection string (Supabase, Railway, etc.)
		dsn = databaseURL

		// For Supabase/Neon: Ensure SSL is required and add connection parameters
		if contains(dsn, "supabase") || contains(dsn, "neon.tech") {
			// Add SSL mode if not present
			if !contains(dsn, "sslmode=") {
				separator := "?"
				if strings.Contains(dsn, "?") {
					separator = "&"
				}
				dsn += separator + "sslmode=require"
			}
			// Force IPv4 and add connection timeout
			if !contains(dsn, "connect_timeout=") {
				separator := "?"
				if strings.Contains(dsn, "?") {
					separator = "&"
				}
				dsn += separator + "connect_timeout=10"
			}
		}

		log.Println("📡 Using DATABASE_URL connection string")
	} else {
		// Fallback to individual environment variables for local development
		host := getEnv("POSTGRES_HOST")
		port := getEnv("POSTGRES_PORT")
		user := getEnv("POSTGRES_USER")
		password := getEnv("POSTGRES_PASSWORD")
		dbname := getEnv("POSTGRES_DB")

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
		log.Println("🔧 Using individual database configuration")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}
	DB = db
	log.Printf("Connection to Postgres successfully")

}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
