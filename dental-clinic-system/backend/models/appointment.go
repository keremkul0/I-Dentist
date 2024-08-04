package models

import "time"

type Appointment struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    PatientID   uint      `json:"patient_id"`
    Patient     Patient   `json:"patient"`
    DoctorID    uint      `json:"doctor_id"`
    Doctor      Doctor    `json:"doctor"`
    ClinicID    uint      `json:"clinic_id"`
    Clinic      Clinic    `json:"clinic"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
}
