package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	NationalID     string `json:"national_id" gorm:"uniqueIndex"`
	Name           string `json:"name"`
	BirthDate      string `json:"birth_date"`
	ContactInfo    string `json:"contact_info"`
	MedicalHistory string `json:"medical_history"`
	ClinicID       uint   `json:"clinic_id"`
	Clinic         Clinic `gorm:"foreignKey:ClinicID"`
}
