package token

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	Token     string    `json:"token" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
}
