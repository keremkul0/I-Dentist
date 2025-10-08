package authMiddleware

import (
	"context"
	"dental-clinic-system/models/claims"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

type JwtService interface {
	GetJwtKey() []byte
}

type AuthMiddleware struct {
	TokenService TokenService
	jwtService   JwtService
}

func NewAuthMiddleware(tokenService TokenService, jwtService JwtService) *AuthMiddleware {
	return &AuthMiddleware{TokenService: tokenService, jwtService: jwtService}
}

func (auth *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Cookie'den token al
		token := c.Cookies("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No token provided",
			})
		}

		// Token blacklist kontrolü
		if auth.TokenService.IsTokenBlacklisted(ctx, token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is blacklisted",
			})
		}

		// JWT token doğrulama
		claims := &claims.Claims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.jwtService.GetJwtKey(), nil
		})

		if err != nil {
			if errors.Is(jwt.ErrSignatureInvalid, err) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token signature",
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if !jwtToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		return c.Next()
	}
}
