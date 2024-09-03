package main

import (
	"net/http"

	"github.com/AJC232/InfinityStream-backend/utils"

	"github.com/gin-gonic/gin"
)

// Healthz is a health check endpoint
func Healthz(c *gin.Context) {
	utils.JSONResponse(c, http.StatusOK, map[string]string{"status": "ok"})
}
