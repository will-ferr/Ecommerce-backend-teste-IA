package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, bool)
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) bool
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
		PoolSize: 10,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: rdb}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Cache get error")
		return nil, false
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Cache unmarshal error")
		return nil, false
	}

	return result, true
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Clear(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) bool {
	_, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Cache exists check error")
		return false
	}
	return true
}
