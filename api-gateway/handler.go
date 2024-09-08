package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz is a health check endpoint
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
