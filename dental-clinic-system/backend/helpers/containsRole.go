package helpers

import (
	"dental-clinic-system/models"
	"strings"
)

func ContainsRole(user models.UserGetModel, roleName string) bool {
	for _, role := range user.Roles {
		if strings.EqualFold(role.Name, roleName) {
			return true
		}
	}
	return false
}
