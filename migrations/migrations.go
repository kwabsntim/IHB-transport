package migrations

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"log"
)

func Migrate() {
	// AutoMigrate is now safe - all fields are nullable or have defaults
	// This is idempotent and production-safe
	log.Println("🚀 Running AutoMigrate...")
	err := database.DB.AutoMigrate(
		&models.DeliveryRequest{},
		&models.StatusLog{},
		&models.EmailLog{},
		&models.Admin{},
	)

	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	log.Println("✅ Migration successful - safe to run anytime!")
}
