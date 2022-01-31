package infrastructure

import "github.com/go-redis/redis/v8"

var redisClient *redis.Client

func GetRedisClient(opts *redis.Options) *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = redis.NewClient(opts)
	return redisClient
}
