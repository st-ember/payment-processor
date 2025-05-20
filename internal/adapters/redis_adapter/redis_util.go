package redisadapter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct{}

func NewRedisUtil() *RedisUtil { return &RedisUtil{} }

func (u *RedisUtil) StoreStringRecord(rdb *redis.Client, ctx context.Context, key, record string) error {
	err := rdb.Set(ctx, key, record, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *RedisUtil) GetStringRecord(rdb *redis.Client, ctx context.Context, key string) (string, error) {
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (u *RedisUtil) DelRecord(rdb *redis.Client, ctx context.Context, key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

type StringRecordStorer interface {
	StoreStringRecord(rdb *redis.Client, ctx context.Context, key, record string) error
}

type StringRecordGetter interface {
	GetStringRecord(rdb *redis.Client, ctx context.Context, key string) (string, error)
}
