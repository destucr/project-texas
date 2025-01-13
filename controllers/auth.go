package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/mail"
	"project-texas/config"
	"project-texas/models"
	"project-texas/utils"
	"strings"
)

// / Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON input
	if r.Body == nil {
		jsonError(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate input
	if len(input.Username) == 0 {
		jsonError(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 8 {
		jsonError(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if !isValidEmail(input.Email) {
		jsonError(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		jsonError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Save user to database
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Println("Error:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			jsonError(w, "Email or username already taken", http.StatusBadRequest)
			return
		}
		jsonError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string `json:"identifier"` // Can be email or username
		Password   string `json:"password"`
	}

	// Parse JSON input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	var err error

	// Search by email or username in one query
	err = config.DB.Select("id", "email", "username", "password").
		Where("email = ? OR username = ?", input.Identifier, input.Identifier).
		First(&user).Error

	// Handle errors
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid email/username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// Check password
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		http.Error(w, "Invalid email/username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, expiresIn, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Respond with token and additional information
	response := map[string]interface{}{
		"message":    "Login successful",
		"token":      token,
		"token_type": "Bearer",
		"expires_in": expiresIn,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
