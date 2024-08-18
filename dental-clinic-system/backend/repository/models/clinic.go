package models

import "gorm.io/gorm"

type Clinic struct {
	gorm.Model
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	GroupID     uint   `json:"group_id"`
	Group       Group  `gorm:"foreignKey:GroupID;references:ID"`
	Email       string `json:"email"`
}
