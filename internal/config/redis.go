package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
