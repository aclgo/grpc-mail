package mail

import (
	"context"
	"time"
)

type MailRepo interface {
	Set(ctx context.Context, email string, value any, ttl time.Duration)
	Get(ctx context.Context, email string) (string, error)
}
