package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"konzek-jun/dto"
	"konzek-jun/globalerror"
	"konzek-jun/loggerx"
	"konzek-jun/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type authHandler struct {
	authService services.AuthService
	jwtService  services.JWTService
	userService services.UserService
}

func NewAuthHandler(authService services.AuthService, jwtService services.JWTService, userService services.UserService) AuthHandler {
	return &authHandler{
		authService: authService,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (c *authHandler) Login(ctx *fiber.Ctx) error {
	loggerx.Info("Login function called")

	var loginRequest dto.LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		log.Println("Request parsing error:", err)
		return ctx.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
			Status: http.StatusBadRequest,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Login",
					Description: "Failed to process request",
				},
			},
		})
	}

	loggerx.Info(fmt.Sprintf("Login request received: %+v", loginRequest.Email))

	if errors := globalerror.Validate(loginRequest); len(errors) > 0 && errors[0].HasError {
		loggerx.Info("Invalid login request")
		return globalerror.HandleValidationErrors(ctx, errors)
	}

	log.Printf("Verifying login request: Email - %s", loginRequest.Email)
	err := c.authService.VerifyCredential(loginRequest.Email, loginRequest.Password)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Login verification error: %s", err.Error()))
		return ctx.Status(http.StatusUnauthorized).JSON(globalerror.ErrorResponse{
			Status: http.StatusUnauthorized,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Login",
					Description: "Failed to login",
				},
			},
		})
	}

	user, _ := c.userService.FindUserByEmail(loginRequest.Email)

	token := c.jwtService.GenerateToken(strconv.FormatInt(user.ID, 10))
	user.Token = token
	return ctx.Status(http.StatusOK).JSON(user)
}

func (c *authHandler) Register(ctx *fiber.Ctx) error {
	loggerx.Info("Register function called")

	var registerRequest dto.RegisterRequest
	if err := ctx.BodyParser(&registerRequest); err != nil {
		log.Println("Request parsing error:", err)
		return ctx.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
			Status: http.StatusBadRequest,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Register",
					Description: "Failed to process request",
				},
			},
		})
	}

	log.Println("Register request received:", registerRequest)

	if errors := globalerror.Validate(registerRequest); len(errors) > 0 && errors[0].HasError {
		loggerx.Info("Invalid register request")
		return globalerror.HandleValidationErrors(ctx, errors)
	}

	log.Printf("Creating new user: Email - %s", registerRequest.Email)
	user, err := c.userService.CreateUser(registerRequest)
	if err != nil {
		loggerx.Error(fmt.Sprintf("User creation error: %s", err.Error()))
		return ctx.Status(http.StatusUnprocessableEntity).JSON(globalerror.ErrorResponse{
			Status: http.StatusUnprocessableEntity,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Login",
					Description: "Failed to create, User already exist",
				},
			},
		})
	}

	token := c.jwtService.GenerateToken(strconv.FormatInt(user.ID, 10))
	user.Token = token
	return ctx.Status(http.StatusCreated).JSON(user)
}
