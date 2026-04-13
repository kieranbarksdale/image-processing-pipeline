package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func CreateRedisClient(connectionString string) (*RedisClient, error) {
	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}
