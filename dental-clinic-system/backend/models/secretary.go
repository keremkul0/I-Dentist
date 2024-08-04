package models

type Secretary struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Name      string `json:"name"`
    ClinicID  uint   `json:"clinic_id"`
    Clinic    Clinic `json:"clinic"`
}
