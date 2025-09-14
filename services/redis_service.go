package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-boilerplate/logger"
	repoif "go-boilerplate/repository/interfaces"
	serviceif "go-boilerplate/services/interfaces"
)

// redisService implements service-level Redis operations using a RedisRepository.
type redisService struct {
	repo repoif.RedisRepository
}

func NewRedisService(repo repoif.RedisRepository) serviceif.RedisService {
	return &redisService{repo: repo}
}

func (s *redisService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		logger.Error(ctx, "RedisService.SetJSON marshal failed", map[string]any{"key": key, "error": err.Error()})
		return err
	}
	if err := s.repo.Set(ctx, key, string(b), expiration); err != nil {
		logger.Error(ctx, "RedisService.SetJSON set failed", map[string]any{"key": key, "error": err.Error()})
		return err
	}
	logger.Debug(ctx, "RedisService.SetJSON success", map[string]any{"key": key})
	return nil
}

func (s *redisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := s.repo.Get(ctx, key)
	if err != nil {
		logger.Debug(ctx, "RedisService.GetJSON get failed", map[string]any{"key": key, "error": err.Error()})
		return err
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		logger.Error(ctx, "RedisService.GetJSON unmarshal failed", map[string]any{"key": key, "error": err.Error()})
		return err
	}
	logger.Debug(ctx, "RedisService.GetJSON success", map[string]any{"key": key})
	return nil
}

func (s *redisService) Delete(ctx context.Context, key string) error {
	if _, err := s.repo.Del(ctx, key); err != nil {
		logger.Warn(ctx, "RedisService.Delete failed", map[string]any{"key": key, "error": err.Error()})
		return err
	}
	logger.Debug(ctx, "RedisService.Delete success", map[string]any{"key": key})
	return nil
}

func (s *redisService) Exists(ctx context.Context, keys ...string) (int64, error) {
	res, err := s.repo.Exists(ctx, keys...)
	if err != nil {
		logger.Warn(ctx, "RedisService.Exists failed", map[string]any{"keys": keys, "error": err.Error()})
		return 0, err
	}
	logger.Debug(ctx, "RedisService.Exists success", map[string]any{"keys": keys, "exists": res})
	return res, nil
}

func (s *redisService) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	ok, err := s.repo.Expire(ctx, key, expiration)
	if err != nil {
		logger.Warn(ctx, "RedisService.Expire failed", map[string]any{"key": key, "error": err.Error()})
		return false, err
	}
	logger.Debug(ctx, "RedisService.Expire success", map[string]any{"key": key, "ok": ok})
	return ok, nil
}

func (s *redisService) Incr(ctx context.Context, key string) (int64, error) {
	v, err := s.repo.Incr(ctx, key)
	if err != nil {
		logger.Warn(ctx, "RedisService.Incr failed", map[string]any{"key": key, "error": err.Error()})
		return 0, err
	}
	logger.Debug(ctx, "RedisService.Incr success", map[string]any{"key": key, "value": v})
	return v, nil
}

// Helpers
func (s *redisService) CacheUserSession(ctx context.Context, userID uint, token string, ttl time.Duration) error {
	key := fmt.Sprintf("user_session:%d", userID)
	if err := s.repo.Set(ctx, key, token, ttl); err != nil {
		logger.Warn(ctx, "RedisService.CacheUserSession failed", map[string]any{"user_id": userID, "error": err.Error()})
		return err
	}
	logger.Debug(ctx, "RedisService.CacheUserSession success", map[string]any{"user_id": userID})
	return nil
}

func (s *redisService) GetUserSession(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user_session:%d", userID)
	val, err := s.repo.Get(ctx, key)
	if err != nil {
		logger.Debug(ctx, "RedisService.GetUserSession miss", map[string]any{"user_id": userID, "error": err.Error()})
		return "", err
	}
	logger.Debug(ctx, "RedisService.GetUserSession hit", map[string]any{"user_id": userID})
	return val, nil
}

func (s *redisService) DeleteUserSession(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user_session:%d", userID)
	if _, err := s.repo.Del(ctx, key); err != nil {
		logger.Warn(ctx, "RedisService.DeleteUserSession failed", map[string]any{"user_id": userID, "error": err.Error()})
		return err
	}
	logger.Debug(ctx, "RedisService.DeleteUserSession success", map[string]any{"user_id": userID})
	return nil
}
