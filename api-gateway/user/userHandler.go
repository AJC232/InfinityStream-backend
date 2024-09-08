package user

import (
	"log"
	"net/http"

	proto "github.com/AJC232/InfinityStream-backend/common/protoc"

	"github.com/AJC232/InfinityStream-backend/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userServiceClient proto.UserServiceClient

// Initialize user service client
func InitializeGrpcClient(domain, port string) {
	conn, err := grpc.NewClient(domain+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	userServiceClient = proto.NewUserServiceClient(conn)
}

func RegisterUser(c *gin.Context) {
	// Create a new userRequest struct
	var userRequest proto.UserRegisterRequest

	// Bind JSON request body to the userRequest struct
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	userResponse, err := userServiceClient.RegisterUser(c, &userRequest)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, userResponse)
}

func LoginUser(c *gin.Context) {
	// Create a new userRequest struct
	var userRequest proto.UserLoginRequest

	// Bind JSON request body to the userRequest struct
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	userResponse, err := userServiceClient.LoginUser(c, &userRequest)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, userResponse)
}

func GetUser(c *gin.Context) {
	// Create a new userRequest struct
	var userRequest proto.GetUserRequest

	// Extract the user ID path parameter
	userRequest.Id = c.Param("userId")

	userResponse, err := userServiceClient.GetUser(c, &userRequest)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, userResponse)
}
