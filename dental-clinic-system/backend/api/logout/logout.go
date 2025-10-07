package logout

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TokenService interface {
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
}

type LogoutController struct {
	tokenService TokenService
}

func NewLogoutController(tokenService TokenService) *LogoutController {
	return &LogoutController{
		tokenService: tokenService,
	}
}

func (h *LogoutController) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No token found",
		})
	}

	err := h.tokenService.AddTokenToBlacklist(ctx, token, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Logout failed",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
