package patient

import (
	"dental-clinic-system/models/clinic"
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	NationalID     string        `json:"national_id" gorm:"uniqueIndex"`
	Name           string        `json:"name"`
	BirthDate      string        `json:"birth_date"`
	ContactInfo    string        `json:"contact_info"`
	MedicalHistory string        `json:"medical_history"`
	ClinicID       uint          `json:"clinic_id"`
	Clinic         clinic.Clinic `gorm:"foreignKey:ClinicID"`
}
