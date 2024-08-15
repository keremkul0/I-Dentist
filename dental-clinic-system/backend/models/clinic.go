package models

type Clinic struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	GroupID     uint   `json:"group_id"`
	Email       string `json:"email"` // Eklendi
}
