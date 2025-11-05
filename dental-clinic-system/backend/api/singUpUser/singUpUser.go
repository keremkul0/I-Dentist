package singUpUser

import (
	"context"
	"dental-clinic-system/models/user"

	"github.com/gofiber/fiber/v2"
)

type SignUpUserService interface {
	SignUpUser(ctx context.Context, user user.User) (string, error)
}

type SignUpUserHandler struct {
	signUpUserService SignUpUserService
}

func NewSignUpUserHandler(signUpUserService SignUpUserService) *SignUpUserHandler {
	return &SignUpUserHandler{
		signUpUserService: signUpUserService,
	}
}

func (s *SignUpUserHandler) SignUpUser(c *fiber.Ctx) error {
	ctx := c.Context()
	var user user.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cacheKey, err := s.signUpUserService.SignUpUser(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(cacheKey)
}
