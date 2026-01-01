package migrations

import (
	"IHB-transport/internal/database"
	"IHB-transport/internal/models"
	"log"
)

func Migrate() {
	// Check if the User table already exists
	if !database.DB.Migrator().HasTable(&models.User{}) {
		// this does everything like creating the tables automatically
		err := database.DB.AutoMigrate(
			//getting all the models
			&models.User{},
			&models.DeliveryRequest{},
			&models.StatusLog{},
			&models.EmailLog{},
		)
		if err != nil {
			log.Fatalf("Migration failed...%v", err)
		}
		log.Println("Migration successful")
	} else {
		log.Println("Tables already exist, skipping migration.")
	}
}
