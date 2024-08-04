package models

type Patient struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `json:"name"`
	BirthDate      string `json:"birth_date"`
	ContactInfo    string `json:"contact_info"`
	MedicalHistory string `json:"medical_history"`
}
