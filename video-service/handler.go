package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AJC232/InfinityStream-backend/config"
	"github.com/AJC232/InfinityStream-backend/utils"
	"github.com/AJC232/InfinityStream-backend/video-service/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB
var s3Client *s3.S3
var awsRegion, awsAccessKey, awsSecretKey, s3BucketName string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	db = config.InitializeDB()

	awsRegion = os.Getenv("AWS_REGION")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	s3BucketName = os.Getenv("S3_BUCKET_NAME")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		log.Fatal("Error creating AWS session:", err)
	}
	s3Client = s3.New(sess)
}

// UploadVideo handles video file uploads and saves them to S3
func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the incoming file (e.g., 100 MB max upload size)
	r.ParseMultipartForm(100 << 20) // 100 MB

	// Retrieve the file from the form data
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	video := models.Video{
		ID:        uuid.New().String(),
		Name:      handler.Filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Size:      handler.Size,
		Type:      handler.Header.Get("Content-Type"),
		S3Url:     fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", s3BucketName, awsRegion, handler.Filename),
	}

	if err := db.Create(&video).Error; err != nil {
		fmt.Println("Error inserting video:", err)
		utils.JSONError(w, http.StatusInternalServerError, "Error inserting video")
		return
	}

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(video.Name),
		Body:        file,
		ContentType: aws.String(handler.Header.Get("Content-Type")),
	})
	if err != nil {
		fmt.Println("Error uploading video to S3:", err)
		utils.JSONError(w, http.StatusInternalServerError, "Error uploading video to S3")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Video uploaded successfully"})
	fmt.Fprintf(w, "Video uploaded successfully: %s\n", video.Name)
}

// StreamVideo streams video content from S3
func StreamVideo(w http.ResponseWriter, r *http.Request) {
	// Extract the video name from the URL query parameter
	videoName := r.URL.Query().Get("video")

	// Validate the video name
	if videoName == "" {
		utils.JSONError(w, http.StatusBadRequest, "Invalid video name")
		return
	}

	// Get the video from S3
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(videoName),
	})
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "Video not found in S3")
		return
	}
	defer result.Body.Close()

	// Set the appropriate content type
	w.Header().Set("Content-Type", *result.ContentType)

	// Stream the video file to the client
	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, "Error streaming video", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Video streamed successfully"})
}

// func GetAllVideos(w http.ResponseWriter, r *http.Request) {

// }

// func DeleteVideo(w http.ResponseWriter, r *http.Request) {

// }
