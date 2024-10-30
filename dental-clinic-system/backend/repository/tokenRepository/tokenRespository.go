package tokenRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
	"time"
)

type TokenRepository interface {
	DeleteExpiredTokensRepo()
	AddTokenToBlacklistRepo(token string, expireTime time.Time)
	IsTokenBlacklistedRepo(token string) bool
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{DB: db}
}

type tokenRepository struct {
	DB *gorm.DB
}

func (r *tokenRepository) DeleteExpiredTokensRepo() {
	_ = r.DB.Unscoped().Where("expire_time < ?", time.Now().Unix()).Delete(&models.ExpiredTokens{}).Error
	return
}

func (r *tokenRepository) AddTokenToBlacklistRepo(token string, expireTime time.Time) {
	r.DB.Create(&models.ExpiredTokens{
		Token:      token,
		ExpireTime: expireTime.Unix(),
	})
	return
}

func (r *tokenRepository) IsTokenBlacklistedRepo(token string) bool {
	var count int64
	r.DB.Model(&models.ExpiredTokens{}).Where("token = ?", token).Count(&count)
	return count > 0
}
