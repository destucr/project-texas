package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"project-texas/config"
	"project-texas/models"
	"project-texas/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Debugging: Print out the loaded environment variables
	log.Printf("DB_USER: %s", os.Getenv("DB_USER"))
	log.Printf("SECRET_KEY: %s", os.Getenv("SECRET_KEY"))

	config.Connect()
	config.DB.AutoMigrate(&models.User{})

	// Register routes
	routes.RegisterRoutes()

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
