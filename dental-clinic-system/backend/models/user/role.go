package user

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name  RoleName `json:"name"`
	Users []*User  `gorm:"many2many:user_roles;"`
}
