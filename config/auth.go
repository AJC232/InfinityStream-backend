package config

import (
	"context"
	"net/http"
	"time"

	"github.com/AJC232/InfinityStream-backend/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	jwtSecretKey = []byte("infinitystream_secret_key")
)

// Claims defines the structure of the JWT claims
type Claims struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(ID uuid.UUID, username string) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	// Create the claims
	claims := &Claims{
		ID:       ID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "infinitystream",
		},
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// AuthMiddleware is a middleware to authenticate requests
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.JSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil {
			utils.JSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Check if the token is valid
		if !token.Valid {
			utils.JSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Set the user ID in the request context
		type contextKey string
		const userIDKey contextKey = "userID"

		ctx := context.WithValue(r.Context(), userIDKey, claims.ID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetUserInfo gets logged in user info
func GetUserInfo(r *http.Request) uuid.UUID {
	type contextKey string
	const userIDKey contextKey = "userID"
	userID := r.Context().Value(userIDKey).(uuid.UUID)
	return userID
}
