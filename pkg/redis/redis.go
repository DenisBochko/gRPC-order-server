package redisClient

import (
	"context"
	"fmt"
	"order-server/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClientCfg struct {
	Host     string `yaml:"REDIS_HOST" env:"REDIS_HOST"`
	Port     string `yaml:"REDIS_PORT" env:"REDIS_PORT" env-default:"6379"`
	Password string `yaml:"REDIS_PASS" env:"REDIS_PASS" env-default:"password"`
	DB       int    `yaml:"REDIS_DB" env:"REDIS_DB" env-default:"0"`
}

func New(ctx context.Context, config RedisClientCfg) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if response := rdb.Ping(ctx); response.Err() != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to redis", zap.Error(response.Err()))
	}

	return rdb
}
