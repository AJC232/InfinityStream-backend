package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc/video"
	"github.com/AJC232/InfinityStream-backend/config"
	"github.com/AJC232/InfinityStream-backend/video-service/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	db                                                  *gorm.DB
	s3Client                                            *s3.S3
	awsRegion, awsAccessKey, awsSecretKey, s3BucketName string
)

func init() {
	err := godotenv.Load("../.env")
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

type Video struct {
	proto.UnimplementedVideoServiceServer
}

func (s *Video) UploadVideo(ctx context.Context, req *proto.UploadVideoRequest) (*proto.UploadVideoResponse, error) {
	// Generate a new UUID for the video
	videoID := uuid.New().String()

	// Define the file paths in S3
	videoFilePath := "videos/" + videoID + ".mp4"            // Customize path and file extension if needed
	coverPhotoFilePath := "cover_photos/" + videoID + ".jpg" // Customize path and file extension if needed

	// Create request object for video pre-signed URL
	videoPutObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(videoFilePath),
	}

	// Generate a pre-signed URL for video upload
	videoReqObj, _ := s3Client.PutObjectRequest(videoPutObjectInput)
	videoReqURL, err := videoReqObj.Presign(15 * time.Minute) // URL valid for 15 minutes
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating video pre-signed URL")
	}

	// Create request object for cover photo pre-signed URL
	coverPhotoPutObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(coverPhotoFilePath),
	}

	// Generate a pre-signed URL for cover photo upload
	coverPhotoReqObj, _ := s3Client.PutObjectRequest(coverPhotoPutObjectInput)
	coverPhotoReqURL, err := coverPhotoReqObj.Presign(15 * time.Minute) // URL valid for 15 minutes
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating cover photo pre-signed URL")
	}

	// Create a new video record in the database
	video := models.Video{
		ID:          videoID,
		Title:       req.Title,
		Description: req.Description,
		Type:        "mp4", // Customize file extension if needed
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Category:    req.Category,
		IsPremium:   req.IsPremium,
	}

	err = db.Create(&video).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "Error creating video record in database")
	}

	// Return the response containing the pre-signed URLs and video ID
	uploadVideoResponse := &proto.UploadVideoResponse{
		VideoSignedUrl: videoReqURL,
		CoverPhotoUrl:  coverPhotoReqURL,
		VideoId:        videoID,
	}

	return uploadVideoResponse, nil
}

func (s *Video) UploadCallback(ctx context.Context, req *proto.UploadCallbackRequest) (*proto.UploadCallbackResponse, error) {
	// Extract video ID from the request
	videoID := req.GetVideoId()
	if videoID == "" {
		return nil, status.Error(codes.InvalidArgument, "Video ID is required")
	}

	// Construct the actual S3 URLs
	actualVideoURL := fmt.Sprintf("https://%v.s3.amazonaws.com/videos/%v.mp4", s3BucketName, videoID)
	actualCoverPhotoURL := fmt.Sprintf("https://%v.s3.amazonaws.com/cover_photos/%v.jpg", s3BucketName, videoID)

	// Update the video metadata in the database
	err := db.Model(&models.Video{}).Where("id = ?", videoID).Updates(models.Video{
		VideoUrl:      actualVideoURL,
		CoverPhotoURL: actualCoverPhotoURL,
	}).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "Error updating video metadata: "+err.Error())
	}

	uploadCallbackResponse := &proto.UploadCallbackResponse{
		Message: "Video metadata updated successfully",
	}
	// Return success response
	return uploadCallbackResponse, nil
}

func (s *Video) GetVideoMetadata(ctx context.Context, req *proto.GetVideoMetadataRequest) (*proto.GetVideoMetadataResponse, error) {
	// Extract video ID from the request
	videoID := req.VideoId

	// Retrieve video metadata from the database
	var video models.Video
	err := db.Where("id = ?", videoID).First(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Video not found")
		}
		return nil, status.Error(codes.Internal, "Error retrieving video metadata: "+err.Error())
	}

	// Get s3Key from the video URL
	videoKey := video.VideoUrl[len("https://"+s3BucketName+".s3.amazonaws.com/"):]

	// Generate pre-signed URL for video and cover photo
	videoURL, err := s.generatePresignedURL(videoKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating pre-signed URL for video")
	}

	// Get s3Key from the cover photo URL
	coverPhotoKey := video.CoverPhotoURL[len("https://"+s3BucketName+".s3.amazonaws.com/"):]

	coverPhotoURL, err := s.generatePresignedURL(coverPhotoKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating pre-signed URL for cover photo")
	}

	// Construct the response
	response := &proto.GetVideoMetadataResponse{
		Id:            video.ID,
		Title:         video.Title,
		Description:   video.Description,
		VideoUrl:      videoURL,
		CoverPhotoUrl: coverPhotoURL,
		Type:          video.Type,
		Category:      video.Category,
		IsPremium:     video.IsPremium,
	}

	return response, nil
}

func (s *Video) ListVideos(ctx context.Context, req *proto.ListVideosRequest) (*proto.ListVideosResponse, error) {
	// Extract parameters from the request
	category := req.Category
	onlyPremium := req.OnlyPremium

	// Initialize the query builder
	var query *gorm.DB
	query = db.Model(&models.Video{})

	// Filter by category if provided
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Filter by premium status if requested
	if onlyPremium {
		query = query.Where("is_premium = ?", true)
	}

	// Execute the query and fetch results
	var videos []models.Video
	err := query.Find(&videos).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "Error retrieving videos: "+err.Error())
	}

	// Map the results to the response format
	var videoMetadataList []*proto.GetVideoMetadataResponse
	for _, video := range videos {
		// Get s3Key from the video URL
		videoKey := video.VideoUrl[len("https://"+s3BucketName+".s3.amazonaws.com/"):]

		// Generate pre-signed URL for video and cover photo
		videoURL, err := s.generatePresignedURL(videoKey)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error generating pre-signed URL for video")
		}

		// Get s3Key from the cover photo URL
		coverPhotoKey := video.CoverPhotoURL[len("https://"+s3BucketName+".s3.amazonaws.com/"):]

		coverPhotoURL, err := s.generatePresignedURL(coverPhotoKey)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error generating pre-signed URL for cover photo")
		}

		videoMetadata := &proto.GetVideoMetadataResponse{
			Id:            video.ID,
			Title:         video.Title,
			Description:   video.Description,
			VideoUrl:      videoURL,
			CoverPhotoUrl: coverPhotoURL,
			Category:      video.Category,
			IsPremium:     video.IsPremium,
		}
		videoMetadataList = append(videoMetadataList, videoMetadata)
	}

	// Construct the response
	response := &proto.ListVideosResponse{
		Videos: videoMetadataList,
	}

	return response, nil
}

// Helper function to generate a pre-signed URL for S3
func (s *Video) generatePresignedURL(filePath string) (string, error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(filePath),
	})
	signedURL, err := req.Presign(15 * time.Minute) // URL valid for 15 minutes
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
