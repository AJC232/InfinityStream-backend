package config

import (
	"fmt"
	"log"
	"os"

	userModels "github.com/AJC232/InfinityStream-backend/user-service/models"
	videoModels "github.com/AJC232/InfinityStream-backend/video-service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDB sets up the database connection and performs migrations
func InitializeDB() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err.Error())
		return nil
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPort, dbPassword)

	// Open a connection to the PostgreSQL database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil
	}

	log.Println("Database connection established")

	// Run migrations to create/update tables based on models
	migrate()

	return DB
}

// migrate automatically creates or updates tables based on models
func migrate() {
	// List all your models here
	err := DB.AutoMigrate(
		&userModels.User{},
		&videoModels.Video{},
	// &OtherModel{}, // Add more models if needed
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	log.Println("Database migrated")
}

func GetDB() *gorm.DB {
	return DB
}
