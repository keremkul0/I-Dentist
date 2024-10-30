package models

import "gorm.io/gorm"

type Clinic struct {
	gorm.Model
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number" gorm:"uniqueIndex"`
	Email       string `json:"email" gorm:"uniqueIndex"`
}
