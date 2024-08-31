package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", UploadVideo).Methods("POST")
	r.HandleFunc("/stream/{videoId}", StreamVideo).Methods("GET")
	// r.HandleFunc("/allVideos", GetAllVideos).Methods("GET")
	// r.HandleFunc("/delete", DeleteVideo).Methods("DELETE")

	log.Println("Video Service running on :8082")
	http.ListenAndServe(":8082", r)
}
