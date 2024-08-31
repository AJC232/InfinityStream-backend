package main

import (
	"log"
	"net/http"

	"github.com/AJC232/infinitystream/infinitystream_backend/config"
	"github.com/gorilla/mux"
)

func main() {

	config.InitializeDB()

	r := mux.NewRouter()

	r.HandleFunc("/register", RegisterUser).Methods("POST")
	r.HandleFunc("/login", LoginUser).Methods("POST")

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(config.AuthMiddleware)
	authRouter.HandleFunc("/allusers", GetAllUsers).Methods("GET")
	authRouter.HandleFunc("/user/{id}", GetUser).Methods("GET")

	log.Println("User Service running on :8081")
	http.ListenAndServe(":8081", r)
}
