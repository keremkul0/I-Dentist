package models

type UserRole struct {
    ID     uint `gorm:"primaryKey" json:"id"`
    UserID uint `json:"user_id"`
    RoleID uint `json:"role_id"`
}
