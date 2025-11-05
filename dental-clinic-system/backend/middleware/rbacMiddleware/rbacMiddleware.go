package rbacMiddleware

import (
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/user"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RequireRole(allowedRoles ...user.RoleName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals("user").(*claims.Claims)
		if !ok {
			log.Warn().Str("operation", "RequireRole").Msg("User claims not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		// Kullanıcının sahip olduğu tüm rolleri kontrol et
		userRoles := make([]user.RoleName, len(userClaims.Roles))
		for i, role := range userClaims.Roles {
			userRoles[i] = role.Name
		}

		// İzin verilen rollerden herhangi birine sahip mi kontrol et
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					log.Info().
						Str("operation", "RequireRole").
						Str("user_email", userClaims.Email).
						Str("matched_role", string(userRole)).
						Str("endpoint", c.Path()).
						Msg("Access granted")
					return c.Next()
				}
			}
		}

		// Yetkisiz erişim - kullanıcının hiçbir rolü eşleşmiyor
		userRoleNames := make([]string, len(userRoles))
		for i, role := range userRoles {
			userRoleNames[i] = string(role)
		}

		log.Warn().
			Str("operation", "RequireRole").
			Str("user_email", userClaims.Email).
			Strs("user_roles", userRoleNames).
			Str("endpoint", c.Path()).
			Str("method", c.Method()).
			Msg("Access denied - insufficient permissions")

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Insufficient permissions",
			"message": "You don't have the required role to access this resource",
		})
	}
}

func RequireSuperAdmin() fiber.Handler {
	return RequireRole(user.RoleSuperAdmin)
}

func RequireClinicAdmin() fiber.Handler {
	return RequireRole(user.RoleClinicAdmin)
}
