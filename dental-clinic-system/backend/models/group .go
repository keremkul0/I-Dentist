package models

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	gorm.Model
	Name         string    `json:"name" gorm:"unique;not null"`
	Description  string    `json:"description"`
	LogoURL      string    `json:"logo_url"`
	FoundedDate  time.Time `json:"founded_date"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Type         string    `json:"type"`
}
