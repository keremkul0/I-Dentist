package appointment

import (
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	gorm.Model
	ClinicID      uint          `json:"clinic_id"`
	Clinic        clinic.Clinic `gorm:"foreignKey:ClinicID"`
	PatientID     uint          `json:"patient_id"`
	Patient       user.User     `gorm:"foreignKey:PatientID"`
	DoctorID      uint          `json:"doctor_id"`
	Doctor        user.User     `gorm:"foreignKey:DoctorID"`
	ScheduledTime time.Time     `json:"scheduled_time"`
	Treatment     string        `json:"treatment"`
	Notes         string        `json:"notes"`
}
