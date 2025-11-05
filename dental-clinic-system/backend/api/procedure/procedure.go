package procedure

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/procedure"
	"dental-clinic-system/models/user"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProcedureService interface {
	GetProcedures(ctx context.Context, ClinicID uint) ([]procedure.Procedure, error)
	GetProcedure(ctx context.Context, id uint) (procedure.Procedure, error)
	CreateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
	UpdateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
	DeleteProcedure(ctx context.Context, id uint) error
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type RoleService interface {
	UserHasRole(user user.UserGetModel, roleName string) bool
}

type JwtService interface {
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

func NewProcedureController(procedureService ProcedureService, userService UserService, roleService RoleService, jwtService JwtService) *ProcedureHandler {
	return &ProcedureHandler{
		procedureService: procedureService,
		userService:      userService,
		roleService:      roleService,
		jwtService:       jwtService,
	}
}

type ProcedureHandler struct {
	procedureService ProcedureService
	userService      UserService
	roleService      RoleService
	jwtService       JwtService
}

func (h *ProcedureHandler) GetProcedures(c *fiber.Ctx) error {
	ctx := c.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	procedures, err := h.procedureService.GetProcedures(ctx, user.ClinicID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(procedures)
}

func (h *ProcedureHandler) GetProcedure(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid procedure ID",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	procedure, err := h.procedureService.GetProcedure(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if procedure.ClinicID != user.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(procedure)
}

func (h *ProcedureHandler) CreateProcedure(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	var procedure procedure.Procedure
	err := c.BodyParser(&procedure)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	procedure.ClinicID = user.ClinicID
	procedure, err = h.procedureService.CreateProcedure(ctx, procedure)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create procedure",
		})
	}
	return c.Status(fiber.StatusOK).JSON(procedure)
}

func (h *ProcedureHandler) UpdateProcedure(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	var procedure procedure.Procedure
	err := c.BodyParser(&procedure)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if !(h.roleService.UserHasRole(user, "Clinic Admin") || h.roleService.UserHasRole(user, "Superadmin")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	procedure.ClinicID = user.ClinicID
	procedure, err = h.procedureService.UpdateProcedure(ctx, procedure)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(procedure)
}

func (h *ProcedureHandler) DeleteProcedure(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid procedure ID",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if !(h.roleService.UserHasRole(user, "Clinic Admin") || h.roleService.UserHasRole(user, "Superadmin")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	procedure, err := h.procedureService.GetProcedure(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if procedure.ClinicID != user.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err = h.procedureService.DeleteProcedure(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Procedure deleted successfully",
	})
}
