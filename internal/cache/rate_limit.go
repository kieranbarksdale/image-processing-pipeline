package cache

import (
	"context"
	"time"
)

func (redisClient *RedisClient) IsRateLimited(ip string) (bool, error) {
	key := "rate_limit:" + ip
	counter, err := redisClient.client.Incr(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if counter == 1 {
		err = redisClient.client.Expire(context.Background(), key, 60*time.Second).Err()
		if err != nil {
			return false, err
		}
	}
	if counter >= 10 {
		return true, nil
	}
	
	return false, nil
}
