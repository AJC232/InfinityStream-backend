package config

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jwtSecretKey = []byte("infinitystream_secret_key")

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
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set the user ID in the Gin context
		c.Set("userID", claims.ID)

		// Continue to the next handler
		c.Next()
	}
}

// GetUserInfo gets logged in user info
func GetUserInfo(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("userID")
	if !exists {
		// Handle the case where the userID doesn't exist in the context
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}
