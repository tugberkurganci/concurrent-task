package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"konzek-jun/dto"
	services "konzek-jun/mocks/service"
)

func TestAuthHandler_Login(t *testing.T) {
	// Mock services
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMockService := services.NewMockAuthService(ctrl)
	jwtMockService := services.NewMockJWTService(ctrl)
	userMockService := services.NewMockUserService(ctrl)

	// Create AuthHandler instance
	authHandler := NewAuthHandler(authMockService, jwtMockService, userMockService)
	router := fiber.New()
	router.Post("/api/login", authHandler.Login)
	// Mock login request
	loginRequest := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	// Mock user data
	mockUser := dto.UserResponse{
		ID:    1,
		Email: loginRequest.Email,

		// Add other necessary fields here
	}

	// Mock AuthService.VerifyCredential to return no error
	authMockService.EXPECT().VerifyCredential(loginRequest.Email, loginRequest.Password).Return(nil)

	// Mock UserService.FindUserByEmail to return the mock user
	userMockService.EXPECT().FindUserByEmail(loginRequest.Email).Return(&mockUser, nil)

	// Mock JWTService.GenerateToken to return a token
	jwtMockService.EXPECT().GenerateToken("1").Return("mock_token")

	// Prepare request
	jsonData, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, _ := router.Test(req)

	// Assert response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert response body
	// Add more assertions as per your response structure
	// Example: assert.Equal(t, expectedResponseBody, resp.Body)
}

func TestAuthHandler_Register(t *testing.T) {
	// Mock services
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMockService := services.NewMockAuthService(ctrl)
	jwtMockService := services.NewMockJWTService(ctrl)
	userMockService := services.NewMockUserService(ctrl)

	// Create AuthHandler instance
	authHandler := NewAuthHandler(authMockService, jwtMockService, userMockService)
	router := fiber.New()
	router.Post("/api/register", authHandler.Register)
	// Mock register request
	registerRequest := dto.RegisterRequest{
		Name:     "tugberk",
		Email:    "test@example.com",
		Password: "password",
		// Add other necessary fields here
	}

	mockUser := dto.UserResponse{
		ID:    1,
		Email: registerRequest.Email,
	}
	// Mock UserService.CreateUser to return no error
	userMockService.EXPECT().CreateUser(registerRequest).Return(&mockUser, nil)

	// Mock JWTService.GenerateToken to return a token
	jwtMockService.EXPECT().GenerateToken(gomock.Any()).Return("mock_token")

	// Prepare request
	jsonData, _ := json.Marshal(registerRequest)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, _ := router.Test(req)

	// Assert response status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Assert response body
	// Add more assertions as per your response structure
	// Example: assert.Equal(t, expectedResponseBody, resp.Body)
}
