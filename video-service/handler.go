package main

import (
	"context"
	"log"
	"os"
	"time"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc/video"
	"github.com/AJC232/InfinityStream-backend/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
	"gorm.io/gorm"
)

var (
	db                                                  *gorm.DB
	s3Client                                            *s3.S3
	awsRegion, awsAccessKey, awsSecretKey, s3BucketName string
)

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

type Video struct {
	proto.UnimplementedVideoServiceServer
}

func (s *Video) UploadVideo(ctx context.Context, req *proto.UploadVideoRequest) (*proto.UploadVideoResponse, error) {
	// Define S3 bucket and file path
	videoID := uuid.New().String()
	filePath := "videos/" + videoID + ".mp4" // Customize path and file extension if needed

	// Create a request object
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(filePath),
	}

	// Generate a pre-signed URL for uploading
	reqObj, _ := s3Client.PutObjectRequest(putObjectInput) // Get the request object
	reqURL, err := reqObj.Presign(15 * time.Minute)        // URL valid for 15 minutes
	if err != nil {
		return nil, status.Error(codes.Internal, "Error generating pre-signed URL")
	}

	uploadVideoResponse := &proto.UploadVideoResponse{
		UploadUrl: reqURL,
		VideoId:   videoID,
	}

	return uploadVideoResponse, nil
}

func (s *Video) GetVideoMetadata(ctx context.Context, req *proto.GetVideoMetadataRequest) (*proto.GetVideoMetadataResponse, error) {
	// Implement metadata retrieval logic

	return nil, nil
}

func (s *Video) StreamVideo(ctx context.Context, req *proto.StreamVideoRequest) (*proto.StreamVideoResponse, error) {
	// Implement video streaming logic

	return nil, nil
}

func (s *Video) ListVideos(ctx context.Context, req *proto.ListVideosRequest) (*proto.ListVideosResponse, error) {
	// Implement video listing logic

	return nil, nil
}
