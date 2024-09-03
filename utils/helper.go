package utils

import (
	"log"

    "github.com/gin-gonic/gin"
)

// JSONError sends a JSON error response with the appropriate status code
func JSONError(c *gin.Context, code int, message string) {
    if code > 499 {
        log.Printf("Error: %v", message)
    }

    JSONResponse(c, code, gin.H{"error": message})
}

// JSONResponse sends a JSON response with the appropriate status code
func JSONResponse(c *gin.Context, code int, payload interface{}) {
    c.JSON(code, payload)
}
