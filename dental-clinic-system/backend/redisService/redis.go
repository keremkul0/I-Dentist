package redisService

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var Rdb *redis.Client
var ctx = context.Background()

func InitializeRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis'e bağlanılamadı: %v", err)
	}
	log.Println("Redis bağlantısı başarılı!")
}
