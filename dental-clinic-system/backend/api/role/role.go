package role

import (
	"context"
	"dental-clinic-system/models/user"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]user.Role, error)
	GetRole(ctx context.Context, id uint) (user.Role, error)
	CreateRole(ctx context.Context, role user.Role) (user.Role, error)
	UpdateRole(ctx context.Context, role user.Role) (user.Role, error)
	DeleteRole(ctx context.Context, id uint) error
}

type RoleHandler struct {
	roleService RoleService
}

func NewRoleController(roleService RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	roles, err := h.roleService.GetRoles(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch roles")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}

	return c.Status(fiber.StatusOK).JSON(roles)
}

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid role ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	role, err := h.roleService.GetRole(ctx, uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Role not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	var role user.Role
	if err := c.BodyParser(&role); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	createdRole, err := h.roleService.CreateRole(ctx, role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create role")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create role",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdRole)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.Context()

	var role user.Role
	if err := c.BodyParser(&role); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updatedRole, err := h.roleService.UpdateRole(ctx, role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update role")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedRole)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid role ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	err = h.roleService.DeleteRole(ctx, uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete role")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
