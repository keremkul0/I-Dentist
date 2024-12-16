package tokenRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
	"time"
)

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{DB: db}
}

type tokenRepository struct {
	DB *gorm.DB
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) {
	_ = r.DB.WithContext(ctx).Unscoped().Where("expire_time < ?", time.Now().Unix()).Delete(&models.ExpiredTokens{}).Error
	return
}

func (r *tokenRepository) AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error {
	return r.DB.WithContext(ctx).Create(&models.ExpiredTokens{
		Token:      token,
		ExpireTime: expireTime.Unix(),
	}).Error
}

func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, token string) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.ExpiredTokens{}).Where("token = ?", token).Count(&count)
	return count > 0
}
