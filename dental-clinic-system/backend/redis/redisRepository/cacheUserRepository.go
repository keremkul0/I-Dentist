package redisRepository

import (
	"context"
	"dental-clinic-system/models"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

func NewRedisUserRepository(rdb *redis.Client) *redisUserRepository {
	return &redisUserRepository{
		rdb: rdb,
	}
}

type redisUserRepository struct {
	rdb *redis.Client
}

func (r *redisUserRepository) CacheUserRepo(user models.User) (string, error) {
	ctx := context.Background()
	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", errors.New("User JSON convert error")
	}

	cacheKey := uuid.New().String()

	if r.rdb == nil {
		return "", errors.New("Redis client is nil")
	}

	err = r.rdb.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
		return "", errors.New("Redis set error")
	}

	return cacheKey, nil
}

func (r *redisUserRepository) GetUserRepo(cacheKey string) (models.User, error) {
	ctx := context.Background()
	userJSON, err := r.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return models.User{}, errors.New("Redis get error")
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return models.User{}, errors.New("User JSON unmarshal error")
	}

	return user, nil
}
