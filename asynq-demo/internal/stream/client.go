package stream

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

// NewRedisClient crea una nueva conexión a Redis
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
