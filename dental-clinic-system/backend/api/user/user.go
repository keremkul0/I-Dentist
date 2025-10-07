package user

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/errors"
	"dental-clinic-system/models/user"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetUsers(ctx context.Context, ClinicID uint) ([]user.UserGetModel, error)
	GetUser(ctx context.Context, id uint) (user.UserGetModel, error)
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
	CreateUser(ctx context.Context, user user.User) (user.UserGetModel, error)
	UpdateUser(ctx context.Context, user user.User) (user.UserGetModel, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, user user.UserGetModel) (bool, error)
	CreateUserWithAuthorization(ctx context.Context, newUser user.User, authUserEmail string) (user.UserGetModel, error)
}

type RoleService interface {
	UserHasRole(user user.UserGetModel, roleName string) bool
}

type JwtService interface {
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

type UserHandler struct {
	userService UserService
	roleService RoleService
	jwtService  JwtService
}

func NewUserController(service UserService, roleService RoleService, jwtService JwtService) *UserHandler {
	return &UserHandler{userService: service, roleService: roleService, jwtService: jwtService}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	users, err := h.userService.GetUsers(ctx, authenticatedUser.ClinicID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid authenticatedUser ID",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if authenticatedUser.ClinicID != requestedUser.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(authenticatedUser)
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(authenticatedUser)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newUser user.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := h.userService.CreateUserWithAuthorization(ctx, newUser, claims.Email)
	if err != nil {
		switch err.(type) {
		case *errors.UnauthorizedError:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case *errors.ValidationError:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	ctx := context.Background()

	var updateUser user.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if authenticatedUser.ClinicID != updateUser.ClinicID ||
		(!h.roleService.UserHasRole(authenticatedUser, "Clinic Admin") &&
			!h.roleService.UserHasRole(authenticatedUser, "Superadmin")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	updatedUser, err := h.userService.UpdateUser(ctx, updateUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := context.Background()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Yetki kontrolü
	if authenticatedUser.ClinicID != requestedUser.ClinicID ||
		(!h.roleService.UserHasRole(authenticatedUser, "Clinic Admin") &&
			!h.roleService.UserHasRole(authenticatedUser, "Superadmin")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if h.roleService.UserHasRole(requestedUser, "Superadmin") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot delete Superadmin"})
	}

	if err := h.userService.DeleteUser(ctx, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
