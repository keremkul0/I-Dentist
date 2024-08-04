package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name     string     `json:"name"`
    Email    string     `json:"email" gorm:"unique"`
    Password string     `json:"-"`
    Roles    []Role     `gorm:"many2many:user_roles;" json:"roles"`
}
