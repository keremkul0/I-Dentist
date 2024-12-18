package init_func

import (
	"dental-clinic-system/config"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func ConnectRedis(RDBConfig config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RDBConfig.Addr,
		Password: RDBConfig.Password,
		DB:       RDBConfig.DB,
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Redis")
		panic(err)
	}
	log.Info().Msg("Redis connection successful")
	return rdb
}
