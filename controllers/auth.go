package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"project-texas/config"
	"project-texas/models"
	"project-texas/utils"

	"gorm.io/gorm"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if len(input.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if !isValidEmail(input.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Save user to database
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		if gorm.ErrDuplicatedKey == err {
			http.Error(w, "Email or username already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond
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

	// Check if input is email or username
	isEmail := isValidEmail(input.Identifier)

	var user models.User
	var err error
	if isEmail {
		// Search by email
		err = config.DB.Where("email = ?", input.Identifier).First(&user).Error
	} else {
		// Search by username
		err = config.DB.Where("username = ?", input.Identifier).First(&user).Error
	}

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

// Helper function to validate email
func isValidEmail(input string) bool {
	// Regex pattern for validating email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(input)
}
