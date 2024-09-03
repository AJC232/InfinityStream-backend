package main

import (
	"log"

	"github.com/AJC232/InfinityStream-backend/api-gateway/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	user.InitializeGrpcClient("localhost", ":8081")
	// video.InitializeGrpcClient("localhost", ":8082")
}

func main() {
	r := gin.Default()

	// Configure CORS settings
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://example.com", "http://localhost:3000"}, // Allow specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},     // Allow specific methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},     // Allow specific headers
		ExposeHeaders:    []string{"Content-Length"},                              // Expose specific headers
		AllowCredentials: true,                                                    // Allow credentials (cookies, etc.)
	}))

	// Health check endpoint
	r.GET("/api/healthz", Healthz)

	// Route to User Service
	r.POST("/api/users/register", user.RegisterUser)
	r.POST("/api/users/login", user.LoginUser)
	r.GET("/api/users/user/:userId", user.GetUser)

	// Route to Video Service
	// r.HandleFunc("/api/videos/{path:.*}", HandleVideoService)

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
