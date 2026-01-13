package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}

func PingRedis(client *redis.Client) error {
	ctx := context.Background()
	return client.Ping(ctx).Err()
}
