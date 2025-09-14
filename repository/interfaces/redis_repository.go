package interfaces

import (
	"context"
	"time"
)

// RedisRepository defines low-level Redis operations used by the app.
// Keep this limited to common primitives so it remains generic.
type RedisRepository interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
}
