package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type mailRepo struct {
	redisClient *redis.Client
}

func NewmailRepo(rdsClient *redis.Client) *mailRepo {
	return &mailRepo{
		redisClient: rdsClient,
	}
}

func (m *mailRepo) Set(ctx context.Context, email string, value any, ttl time.Duration) error {
	return nil
}

func (m *mailRepo) Get(ctx context.Context, email string) (string, error) {
	return "", nil
}
