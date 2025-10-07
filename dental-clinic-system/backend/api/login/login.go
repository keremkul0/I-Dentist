package login

import (
	"context"
	"dental-clinic-system/models/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginService interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
}

type JwtService interface {
	GenerateJWTToken(email string, time time.Time) (string, error)
}

type LoginHandler struct {
	loginService LoginService
	jwtService   JwtService
}

func NewLoginController(service LoginService, jwtService JwtService) *LoginHandler {
	return &LoginHandler{loginService: service, jwtService: jwtService}
}

func (h *LoginHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	var creds auth.Login
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	user, err := h.loginService.Login(ctx, creds.Email, creds.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	tokenString, err := h.jwtService.GenerateJWTToken(user.Email, expirationTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}
