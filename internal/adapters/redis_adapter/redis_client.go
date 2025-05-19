package redisadapter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClientFactory struct{}

func NewRedisClientFactory() *RedisClientFactory { return &RedisClientFactory{} }

func (f *RedisClientFactory) InitRedisClient(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

type RedisClientInitter interface {
	InitRedisClient(ctx context.Context) *redis.Client
}
