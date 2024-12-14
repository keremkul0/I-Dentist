package redisRepository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

func NewRedisRepository(rdb *redis.Client) *redisRepository {
	return &redisRepository{
		rdb: rdb,
	}
}

type redisRepository struct {
	rdb *redis.Client
}

func (r *redisRepository) SetData(data any) (string, error) {

	if data == nil {
		return "", errors.New("data is nil")
	}

	marshalledData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return "", errors.New("marshal error")
	}

	ctx := context.Background()
	cacheKey := uuid.New().String()
	err = r.rdb.Set(ctx, cacheKey, marshalledData, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
		return "", errors.New("redis set error")
	}

	return cacheKey, nil
}

func (r *redisRepository) GetData(cacheKey string) (any, error) {
	ctx := context.Background()
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Failed to get cache: %v", err)
		return nil, errors.New("redis get error")
	}

	var unmarshalledData any
	err = json.Unmarshal([]byte(data), &unmarshalledData)
	if err != nil {
		log.Printf("Failed to unmarshal data: %v", err)
		return nil, errors.New("unmarshal error")
	}

	return unmarshalledData, nil
}
