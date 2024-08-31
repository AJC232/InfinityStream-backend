package main

import (
	"net/http"

	"github.com/AJC232/InfinityStream-backend/utils"

	"github.com/gorilla/mux"
)

// Healthz is a health check endpoint
func Healthz(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HanleService redirects the request to the appropriate service
func HanleService(w http.ResponseWriter, r *http.Request, service string) {
	newHost := "http://localhost:"
	if service == "users" {
		newHost += "8081/"
	} else if service == "videos" {
		newHost += "8082/"
	} else {
		utils.JSONError(w, http.StatusBadRequest, "Invalid service")
		return
	}

	newPath := mux.Vars(r)["path"]

	// fmt.Println(newHost + newPath)
	http.Redirect(w, r, newHost+newPath, http.StatusSeeOther)
}

// HandleUserService redirects the request to the user service
func HandleUserService(w http.ResponseWriter, r *http.Request) {
	HanleService(w, r, "users")
}

// HandleVideoService redirects the request to the video service
func HandleVideoService(w http.ResponseWriter, r *http.Request) {
	HanleService(w, r, "videos")
}
