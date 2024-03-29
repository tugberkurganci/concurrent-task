package middleware

import (
	"log"
	"net/http"

	"konzek-jun/globalerror"
	"konzek-jun/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type JWTMiddleware struct {
	jwtService services.JWTService
}

func NewJWTMiddleware(jwtService services.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

func (m *JWTMiddleware) AuthorizeJWT(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {

		return c.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
			Status: http.StatusBadRequest,
			ErrorDetail: []globalerror.ErrorResponseDetail{
				{
					FieldName:   "Failed to process request",
					Description: "No token provided",
				},
			},
		})
	}

	token := m.jwtService.ValidateToken(authHeader)
	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		log.Println("Claim[user_id]: ", claims["user_id"])
		log.Println("Claim[issuer] :", claims["issuer"])
		return c.Next()
	}

	return c.Status(http.StatusBadRequest).JSON(globalerror.ErrorResponse{
		Status: http.StatusBadRequest,
		ErrorDetail: []globalerror.ErrorResponseDetail{
			{
				FieldName:   "Error",
				Description: "Your token is not valid",
			},
		},
	})
}
