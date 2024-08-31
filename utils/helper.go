package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONError sends a JSON error response with the appropriate status code
func JSONError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Error: %v", message)
	}

	JSONResponse(w, code, map[string]string{"error": message})
}

// JSONResponse sends a JSON response with the appropriate status code
func JSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON Response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
