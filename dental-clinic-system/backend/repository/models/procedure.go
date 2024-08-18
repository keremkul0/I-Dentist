package models

import (
	"gorm.io/gorm"
)

type Procedure struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	ClinicID    uint   `json:"clinic_id"`
	Clinic      Clinic `gorm:"foreignKey:ClinicID"`
}
