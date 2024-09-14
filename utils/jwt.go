package utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	userIDKey  = "userId"
	authHeader = "Authorization"
)

// Middleware to check JWT and extract userId
func WithJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Expecting userId as a float64 from claims (since JWT uses float64 for numbers)
		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId := int64(userIdFloat) // Convert float64 to int

		// Add userId to context
		ctx := context.WithValue(r.Context(), userIDKey, userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Extract JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get(authHeader)
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

// Parse and validate the JWT token
func parseToken(tokenString string) (jwt.MapClaims, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}

	// Replace "your-secret-key" with your actual JWT secret key
	secretKey := []byte(config.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// CreateToken generates a JWT token for a userId
func CreateToken(userId int64) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}

	// Replace with your own secret key
	secretKey := []byte(config.JWTSecret)

	// Set token expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.JWTExpirationHours))

	// Create JWT claims, including userId and expiration
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime.Unix(),
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Retrieve userId from request context
func GetUserIdFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid claims")
)
