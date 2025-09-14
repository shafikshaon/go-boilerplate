package interfaces

import (
	"context"
	"time"
)

// RedisService defines higher-level Redis operations used by services/handlers.
// It exposes generic primitives plus some typed helpers.
type RedisService interface {
	// Primitives (JSON marshalling can be handled by callers or by specific helpers)
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)

	// Typed helpers
	CacheUserSession(ctx context.Context, userID uint, token string, ttl time.Duration) error
	GetUserSession(ctx context.Context, userID uint) (string, error)
	DeleteUserSession(ctx context.Context, userID uint) error
}
