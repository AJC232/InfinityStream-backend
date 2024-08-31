package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/api/healthz", Healthz).Methods("GET")

	// Route to User Service
	r.HandleFunc("/api/users/{path:.*}", HandleUserService)

	// Route to Video Service
	r.HandleFunc("/api/videos/{path:.*}", HandleVideoService)

	// Define allowed CORS options
	corsOptions := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	allowCredentials := handlers.AllowCredentials()

	log.Println("API Gateway running on :8080")
	http.ListenAndServe(":8080", handlers.CORS(corsOptions, allowedMethods, allowedHeaders, allowCredentials)(r))
}
