package migrations

import (
	"ihb-transport/internal/database"
	"ihb-transport/internal/models"
	"log"
)

func Migrate() {
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.DeliveryRequest{},
		&models.StatusLog{},
		&models.EmailLog{},
	)

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Migration successful")
}
