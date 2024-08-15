package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username  string    `json:"username"`
    Password  string    `json:"password"`
    Role      string    `json:"role"`
    GroupID   uint      `json:"group_id"`
    ClinicID  uint      `json:"clinic_id"`  // Clinic'e referans için eklendi
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    LastLogin time.Time `json:"last_login"`
    IsActive  bool      `json:"is_active"`
}

