package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"project-texas/config"
	"project-texas/models"
	"project-texas/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	config.Connect()
	config.DB.AutoMigrate(&models.User{})

	// Register routes
	routes.RegisterRoutes()

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
