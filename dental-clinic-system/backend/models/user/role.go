package user

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name  string  `json:"name"`
	Users []*User `gorm:"many2many:user_roles;"`
}
