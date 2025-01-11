package routes

import (
    "net/http"
    "project-texas/controllers"
)

func RegisterRoutes() {
	// Public routes
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
}
