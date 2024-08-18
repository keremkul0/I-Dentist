package models

import (
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	gorm.Model
	ClinicID      uint      `json:"clinic_id"`
	Clinic        Clinic    `gorm:"foreignKey:ClinicID"`
	PatientID     uint      `json:"patient_id"`
	Patient       User      `gorm:"foreignKey:PatientID"`
	DoctorID      uint      `json:"doctor_id"`
	Doctor        User      `gorm:"foreignKey:DoctorID"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Treatment     string    `json:"treatment"`
	Notes         string    `json:"notes"`
}
