package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitService interface {
	Allow(ctx context.Context, userID uint, ip string) bool
	GetRemaining(ctx context.Context, userID uint, ip string) (int, time.Duration)
	ResetUserLimits(ctx context.Context, userID uint) error
}

type RedisRateLimit struct {
	client *redis.Client
	mu     sync.Mutex
}

type UserRateLimit struct {
	Requests  int       `json:"requests"`
	Window    int       `json:"window"`
	LastReset time.Time `json:"last_reset"`
}

func NewRedisRateLimit(addr, password string) (*RedisRateLimit, error) {
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

	return &RedisRateLimit{client: rdb}, nil
}

func (r *RedisRateLimit) Allow(ctx context.Context, userID uint, ip string) bool {
	// Check user-specific rate limit
	userKey := fmt.Sprintf("rate_limit:user:%d", userID)
	userAllowed := r.checkUserLimit(ctx, userKey, 100, time.Hour) // 100 requests per hour per user

	// Check IP-based rate limit
	ipKey := fmt.Sprintf("rate_limit:ip:%s", ip)
	ipAllowed := r.checkIPLimit(ctx, ipKey, 20, time.Minute) // 20 requests per minute per IP

	return userAllowed && ipAllowed
}

func (r *RedisRateLimit) checkUserLimit(ctx context.Context, key string, limit int, window time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return true // Allow on Redis errors
	}

	var userLimit UserRateLimit
	if err := json.Unmarshal([]byte(val), &userLimit); err != nil {
		userLimit = UserRateLimit{Window: int(window.Seconds())}
	} else {
		// Reset if window expired
		if time.Since(userLimit.LastReset) > window {
			userLimit.Requests = 0
			userLimit.LastReset = time.Now()
		}
	}

	if userLimit.Requests >= limit {
		return false
	}

	userLimit.Requests++

	// Save back to Redis
	jsonData, _ := json.Marshal(userLimit)
	r.client.Set(ctx, key, jsonData, window)

	return true
}

func (r *RedisRateLimit) checkIPLimit(ctx context.Context, key string, limit int, window time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return true // Allow on Redis errors
	}

	var ipLimit struct {
		Requests int `json:"requests"`
		Window   int `json:"window"`
	}

	if err := json.Unmarshal([]byte(val), &ipLimit); err != nil {
		ipLimit = struct {
			Requests int `json:"requests"`
			Window   int `json:"window"`
		}{Window: int(window.Seconds())}
	}

	if ipLimit.Requests >= limit {
		return false
	}

	ipLimit.Requests++

	// Save back to Redis
	jsonData, _ := json.Marshal(ipLimit)
	r.client.Set(ctx, key, jsonData, window)

	return true
}

func (r *RedisRateLimit) GetRemaining(ctx context.Context, userID uint, ip string) (int, time.Duration) {
	userKey := fmt.Sprintf("rate_limit:user:%d", userID)
	val, err := r.client.Get(ctx, userKey).Result()
	if err != nil {
		return 100, time.Hour
	}

	var userLimit UserRateLimit
	if err := json.Unmarshal([]byte(val), &userLimit); err != nil {
		return 100, time.Hour
	}

	remaining := 100 - userLimit.Requests
	if remaining < 0 {
		remaining = 0
	}

	return remaining, time.Until(userLimit.LastReset.Add(time.Hour))
}

func (r *RedisRateLimit) ResetUserLimits(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("rate_limit:user:%d", userID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
