package tokenRepository

import (
	"dental-clinic-system/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.ExpiredTokens{})
	if err != nil {
		return nil
	}
	return db
}

func TestDeleteExpiredTokensRepo(t *testing.T) {
	db := setupTestDB()
	repo := NewTokenRepository(db)

	// Add an expired token
	expiredToken := models.ExpiredTokens{
		Token:      "expired_token",
		ExpireTime: time.Now().Add(-time.Hour).Unix(),
	}
	db.Create(&expiredToken)

	// Call the method
	repo.DeleteExpiredTokensRepo()

	// Check if the token is deleted
	var count int64
	db.Model(&models.ExpiredTokens{}).Where("token = ?", "expired_token").Count(&count)
	if count != 0 {
		t.Errorf("expected 0 expired tokens, got %d", count)
	}
}

func TestAddTokenToBlacklistRepo(t *testing.T) {
	db := setupTestDB()
	repo := NewTokenRepository(db)

	// Add a token to the blacklist
	token := "blacklisted_token"
	expireTime := time.Now().Add(time.Hour)
	repo.AddTokenToBlacklistRepo(token, expireTime)

	// Check if the token is added
	var count int64
	db.Model(&models.ExpiredTokens{}).Where("token = ?", token).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 blacklisted token, got %d", count)
	}
}

func TestIsTokenBlacklistedRepo(t *testing.T) {
	db := setupTestDB()
	repo := NewTokenRepository(db)

	// Add a token to the blacklist
	token := "blacklisted_token"
	expireTime := time.Now().Add(time.Hour)
	db.Create(&models.ExpiredTokens{
		Token:      token,
		ExpireTime: expireTime.Unix(),
	})

	// Check if the token is blacklisted
	isBlacklisted := repo.IsTokenBlacklistedRepo(token)
	if !isBlacklisted {
		t.Errorf("expected token to be blacklisted")
	}

	// Check if a non-blacklisted token is not blacklisted
	isBlacklisted = repo.IsTokenBlacklistedRepo("non_blacklisted_token")
	if isBlacklisted {
		t.Errorf("expected token to not be blacklisted")
	}
}
