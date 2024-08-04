package models

type Doctor struct {
    ID           uint   `gorm:"primaryKey" json:"id"`
    Name         string `json:"name"`
    Speciality   string `json:"speciality"`
    ClinicID     uint   `json:"clinic_id"`
}
