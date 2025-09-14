package repository

import (
	"context"
	"time"

	repoif "go-boilerplate/repository/interfaces"

	"github.com/redis/go-redis/v9"
)

// redisRepository is a thin wrapper around go-redis client implementing RedisRepository.
type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) repoif.RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepository) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

func (r *redisRepository) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

func (r *redisRepository) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}
