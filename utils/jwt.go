package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

// JWT secret key loaded from environment variable
var jwtSecret = []byte(getJWTSecretKey())

func getJWTSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY is not set in the environment variables")
	}
	log.Printf("Secret key loaded: %s", secret[:5]) 
	return secret
}

// Claims represents the payload of the JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT for the given user ID
func GenerateJWT(userID uint, email string) (string, int64, error) {
    expiresIn := int64(24 * 60 * 60) // 24 hours in seconds
    
    // Set claims
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    // Create a new token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret key
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", 0, err
    }

    return signedToken, expiresIn, nil
}

// ValidateJWT validates the given JWT and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
