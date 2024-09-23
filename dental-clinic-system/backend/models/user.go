package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	NationalID    string    `json:"national_id"`
	Password      string    `json:"password"`
	ClinicID      uint      `json:"clinic_id"`
	Clinic        Clinic    `gorm:"foreignKey:ClinicID"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	LastLogin     time.Time `json:"last_login"`
	IsActive      bool      `json:"is_active"`
	PhoneNumber   string    `json:"phone_number"`
	PhoneVerified bool      `json:"phone_verified"`
	Roles         []*Role   `gorm:"many2many:user_roles;"`
}
