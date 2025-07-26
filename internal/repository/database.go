package repository

import (
	"log"
	"palu-wiki-backend/internal/models" // Import models package

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("palu_wiki.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate will create/update tables based on the defined models
	err = DB.AutoMigrate(&models.Guide{}, &models.UserQuery{}, &models.OfficialUpdate{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	log.Println("Database connection and migration successful.")
}
