package migrations

import (
	"log"
	"project-texas/config"
	"project-texas/models"
)

func Migrate() {
	// Automatically create or migrate the `users` table
	err := config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
}
