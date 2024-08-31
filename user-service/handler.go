package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AJC232/InfinityStream-backend/config"
	"github.com/AJC232/InfinityStream-backend/user-service/models"
	"github.com/AJC232/InfinityStream-backend/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	config.InitializeDB()
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Get database connection
	db := config.GetDB()
	if db == nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error connecting to database")
		return
	}

	// Check if user already exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		utils.JSONError(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Generate UUID for user
	user.ID = uuid.New()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	tnx := db.Create(&user)
	if tnx.Error != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error registering user")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"ID":      user.ID.String(),
		"message": "User registered successfully",
	})
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Get database connection
	db := config.GetDB()
	if db == nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error connecting to database")
		return
	}

	// Check if user exists
	var dbUser models.User
	db.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID == uuid.Nil {
		utils.JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	// Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := config.GenerateToken(dbUser.ID, dbUser.Username)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "Login successful",
		"ID":      dbUser.ID.String(),
		"token":   token,
	})
}

// GetUser handles fetching user profile
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get database connection
	db := config.GetDB()
	if db == nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error connecting to database")
		return
	}

	fmt.Println(userID)
	// Fetch the user profile
	var user models.User
	db.Where("id = ?", userID).Find(&user)

	fmt.Println(user)
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		utils.JSONError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	var userResponse []models.UserResponse
	err = json.Unmarshal(userJSON, &userResponse)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	utils.JSONResponse(w, http.StatusOK, userResponse)
}

// GetAllUsers handles fetching all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get database connection
	db := config.GetDB()
	if db == nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error connecting to database")
		return
	}

	// Fetch all users
	var users []models.User
	db.Find(&users)

	usersJSON, err := json.Marshal(users)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	var userResponses []models.UserResponse
	err = json.Unmarshal(usersJSON, &userResponses)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	utils.JSONResponse(w, http.StatusOK, userResponses)
}
