package main

import (
	"context"
	"log"
	"time"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc/user"
	"github.com/AJC232/InfinityStream-backend/config"
	"github.com/AJC232/InfinityStream-backend/user-service/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	config.InitializeDB()
}

type User struct {
	proto.UnimplementedUserServiceServer
}

// RegisterUser handles user registration
func (u *User) RegisterUser(c context.Context, req *proto.UserRegisterRequest) (*proto.UserRegisterResponse, error) {
	// Get database connection
	db := config.GetDB()
	if db == nil {
		return nil, status.Error(codes.Internal, "Error connecting to database")
	}

	// Check if user already exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error hashing password")
	}

	user := models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tnx := db.Create(&user)
	if tnx.Error != nil {
		return nil, status.Error(codes.Internal, "Error registering user")
	}

	userResponse := &proto.UserRegisterResponse{
		Id:      user.ID.String(),
		Message: "User registered successfully",
	}

	return userResponse, nil
}

// LoginUser handles user login
func (u *User) LoginUser(c context.Context, req *proto.UserLoginRequest) (*proto.UserLoginResponse, error) {
	// Get database connection
	db := config.GetDB()
	if db == nil {
		return nil, status.Error(codes.Internal, "Error connecting to database")
	}

	// Check if user exists
	var user models.User
	db.Where("username = ?", req.Username).First(&user)
	if user.ID == uuid.Nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Generate JWT token
	token, err := config.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating token")
	}

	userResponse := &proto.UserLoginResponse{
		Id:      user.ID.String(),
		Token:   token,
		Message: "Login successful",
	}

	return userResponse, nil
}

// GetUser handles fetching user profile
func (u *User) GetUser(c context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Get database connection
	db := config.GetDB()
	if db == nil {
		return nil, status.Error(codes.Internal, "Error connecting to database")
	}

	// Fetch the user profile
	var user models.User
	db.Where("id = ?", userID).Find(&user)
	if user.ID == uuid.Nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	userResponse := &proto.GetUserResponse{
		Id:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	return userResponse, nil
}
