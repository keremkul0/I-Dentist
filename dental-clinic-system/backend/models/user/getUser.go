package user

import (
	"time"

	"gorm.io/gorm"
)

type UserGetModel struct {
	gorm.Model
	NationalID  string    `json:"national_id"`
	ClinicID    uint      `json:"clinic_id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	LastLogin   time.Time `json:"last_login"`
	IsActive    bool      `json:"is_active"`
	PhoneNumber string    `json:"phone_number"`
	Roles       []*Role   `gorm:"many2many:user_roles;"`
}
