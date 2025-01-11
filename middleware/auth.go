package middleware

import (
	"net/http"
	"strings"
	"project-texas/utils"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set user ID and email in request context (optional)
		r.Header.Set("X-User-ID", string(rune(claims.UserID)))
		r.Header.Set("X-User-Email", claims.Email)

		// Call next handler
		next.ServeHTTP(w, r)
	})
}
