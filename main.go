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
    // Load environment variables
    if err := godotenv.Load(".env"); err != nil {
        log.Fatalf("Error loading .env file: %s", err)
    }

    // Connect to the database
    config.Connect()

    // Run migrations
    config.DB.AutoMigrate(&models.User{})

    // Register routes
    router := routes.RegisterRoutes()

    // Start the server using the router
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
