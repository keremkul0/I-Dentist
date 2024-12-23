package models

import (
	"gorm.io/gorm"
)

type ExpiredTokens struct {
	gorm.Model
	Token      string `json:"token" gorm:"uniqueIndex"`
	ExpireTime int64  `json:"expire_time"`
}
