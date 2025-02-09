package helpers

import (
	"dental-clinic-system/models/user"
	"strings"
)

func ContainsRole(user user.UserGetModel, roleName string) bool {
	for _, role := range user.Roles {
		if strings.EqualFold(role.Name, roleName) {
			return true
		}
	}
	return false
}
