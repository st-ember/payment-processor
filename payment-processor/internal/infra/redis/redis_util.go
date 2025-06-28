package redisadapter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	rdb *redis.Client
}

func NewRedisUtil(rdb *redis.Client) *RedisUtil {
	return &RedisUtil{
		rdb: rdb,
	}
}

func (u *RedisUtil) SetStripeSession(ctx context.Context, sessionId, orderId string) error {
	key := fmt.Sprintf("stripe_checkout:%s", sessionId)
	return u.rdb.Set(ctx, key, orderId, time.Hour*25).Err()
}

func (u *RedisUtil) GetOrderIdFromSession(ctx context.Context, sessionId string) (string, error) {
	key := fmt.Sprintf("stripe_checkout:%s", sessionId)

	val, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
