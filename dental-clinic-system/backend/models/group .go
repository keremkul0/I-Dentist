package models

import "time"

type Group struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"unique;not null"`
	Description  string
	LogoURL      string
	FoundedDate  time.Time
	ContactEmail string
	ContactPhone string
	Type         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
