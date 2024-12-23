package redisRepository

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (r *redisRepository) SetData(ctx context.Context, data any) (string, error) {
	if data == nil {
		return "", errors.New("data is nil")
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return "", errors.New("gob encode error")
	}

	cacheKey := uuid.New().String()
	err = r.rdb.Set(ctx, cacheKey, buf.Bytes(), 5*time.Minute).Err()
	if err != nil {
		return "", errors.New("redis set error")
	}

	return cacheKey, nil
}

func (r *redisRepository) GetData(ctx context.Context, cacheKey string, target any) error {
	data, err := r.rdb.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return errors.New("redis get error")
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(target)
	if err != nil {
		return errors.New("gob decode error")
	}

	return nil
}

func (r *redisRepository) DeleteData(ctx context.Context, ID string) error {
	err := r.rdb.Del(ctx, ID).Err()
	if err != nil {
		return errors.New("redis delete error")
	}

	return nil
}
