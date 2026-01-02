package seed

import(
	"log"
	"ihb-transport/internal/models"
	"gorm.io/gorm"
	"ihb-transport/internal/auth"
)


func DefaultAdminAccount(db *gorm.DB){
	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count>0{
		log.Println("Admin account exists, skipping seeding.")
		return
	}
}

//hashing password 
hashedPassword,err:=auth.