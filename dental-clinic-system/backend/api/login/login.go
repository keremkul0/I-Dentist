package login

import (
	"context"
	"dental-clinic-system/models/auth"
	"dental-clinic-system/models/user"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginService interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
}

type JwtService interface {
	GenerateJWTToken(email string, roles []*user.Role, time time.Time) (string, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type LoginHandler struct {
	loginService LoginService
	jwtService   JwtService
	userService  UserService
}

func NewLoginController(service LoginService, jwtService JwtService, userService UserService) *LoginHandler {
	return &LoginHandler{loginService: service, jwtService: jwtService, userService: userService}
}

func (h *LoginHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	var creds auth.Login
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	authUser, err := h.loginService.Login(ctx, creds.Email, creds.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	user, err := h.userService.GetUserByEmail(ctx, authUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user information",
		})
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	tokenString, err := h.jwtService.GenerateJWTToken(user.Email, user.Roles, expirationTime)
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
