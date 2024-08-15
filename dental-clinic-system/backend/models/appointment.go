package models

import (
    "time"
)

type Appointment struct {
    ID            uint      `gorm:"primaryKey"`
    ClinicID      uint      `json:"clinic_id"`
    Clinic        Clinic    `gorm:"foreignKey:ClinicID"`
    PatientID     uint      `json:"patient_id"`
    Patient       User      `gorm:"foreignKey:PatientID"`
    DoctorID      uint      `json:"doctor_id"`
    Doctor        User      `gorm:"foreignKey:DoctorID"`
    ScheduledTime time.Time `json:"scheduled_time"`
    Treatment     string    `json:"treatment"`
    Notes         string    `json:"notes"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
