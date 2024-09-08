package video

import (
	"log"
	"net/http"
	"strconv"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc"

	"github.com/AJC232/InfinityStream-backend/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var videoServiceClient proto.VideoServiceClient

// Initialize user service client
func InitializeGrpcClient(domain, port string) {
	conn, err := grpc.NewClient(domain+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	videoServiceClient = proto.NewVideoServiceClient(conn)
}

func UploadVideo(c *gin.Context) {
	// Create a new videoRequest struct
	var req proto.UploadVideoRequest

	// Bind JSON request body to the videoRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	res, err := videoServiceClient.UploadVideo(c, &req)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, res)
}

func UploadVideoCallback(c *gin.Context) {
	// Create a new videoRequest struct
	var req proto.UploadCallbackRequest

	// Bind JSON request body to the videoRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	res, err := videoServiceClient.UploadCallback(c, &req)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, res)
}

func GetVideoMetadata(c *gin.Context) {
	// Create a new videoRequest struct
	var req proto.GetVideoMetadataRequest

	req.VideoId = c.Param("videoId")

	res, err := videoServiceClient.GetVideoMetadata(c, &req)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, res)
}

func ListVideos(c *gin.Context) {
	// Create a new videoRequest struct
	var req proto.ListVideosRequest

	req.Category = c.Query("category")
	onlyPremiumStr := c.Query("onlyPremium")

	if onlyPremiumStr != "" {
		var err error
		req.OnlyPremium, err = strconv.ParseBool(onlyPremiumStr)
		if err != nil {
			log.Printf("Error parsing 'onlyPremium' parameter: %v", err)
			utils.JSONError(c, http.StatusBadRequest, "Invalid request")
		}
	}

	res, err := videoServiceClient.ListVideos(c, &req)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, res)
}
