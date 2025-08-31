package procedure

import (
	"dental-clinic-system/models/clinic"

	"gorm.io/gorm"
)

type Procedure struct {
	gorm.Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ClinicID    uint          `json:"clinic_id"`
	Clinic      clinic.Clinic `gorm:"foreignKey:ClinicID"`
}
