package redis

import (
	"context"
	"log"

	"github.com/aclgo/grpc-mail/config"
	"github.com/go-redis/redis/v8"
)

func Connect(config *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Username: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis.Client.Ping: %v", err)
	}

	return client
}
