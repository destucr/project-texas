package routes

import (
    "github.com/gorilla/mux"
    "project-texas/controllers"
)

func RegisterRoutes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/register", controllers.Register).Methods("POST")
    r.HandleFunc("/login", controllers.Login).Methods("POST")
    return r
}