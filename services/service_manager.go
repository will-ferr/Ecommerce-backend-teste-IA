package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type ServiceManager struct {
	cacheService     CacheService
	jobQueue         JobQueue
	rateLimitService RateLimitService
	mu               sync.RWMutex
	initialized      bool
}

var (
	serviceManagerInstance *ServiceManager
	serviceManagerOnce     sync.Once
)

func GetServiceManager() *ServiceManager {
	serviceManagerOnce.Do(func() {
		serviceManagerInstance = &ServiceManager{}
	})
	return serviceManagerInstance
}

func (sm *ServiceManager) InitializeServices() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.initialized {
		return nil
	}

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// Initialize Cache Service
	cache, err := NewRedisCache(redisAddr, redisPassword)
	if err != nil {
		return fmt.Errorf("failed to initialize cache service: %w", err)
	}
	sm.cacheService = cache

	// Initialize Job Queue
	jobQueue, err := NewRedisJobQueue(redisAddr, redisPassword)
	if err != nil {
		return fmt.Errorf("failed to initialize job queue: %w", err)
	}
	sm.jobQueue = jobQueue

	// Initialize Rate Limit Service
	rateLimit, err := NewRedisRateLimit(redisAddr, redisPassword)
	if err != nil {
		return fmt.Errorf("failed to initialize rate limit service: %w", err)
	}
	sm.rateLimitService = rateLimit

	sm.initialized = true
	log.Info().Msg("All services initialized successfully")
	return nil
}

func (sm *ServiceManager) GetCacheService() CacheService {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.cacheService
}

func (sm *ServiceManager) GetJobQueue() JobQueue {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.jobQueue
}

func (sm *ServiceManager) GetRateLimitService() RateLimitService {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.rateLimitService
}

func (sm *ServiceManager) Shutdown(ctx context.Context) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.initialized {
		return nil
	}

	log.Info().Msg("Shutting down services...")

	// Shutdown services gracefully
	var errors []error

	// Note: Redis connections are handled by the client library
	// We just mark services as uninitialized
	sm.initialized = false

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	log.Info().Msg("All services shut down successfully")
	return nil
}

func (sm *ServiceManager) HealthCheck(ctx context.Context) map[string]interface{} {
	health := make(map[string]interface{})

	// Check Cache Service
	if sm.cacheService != nil {
		cacheHealth := "healthy"
		if _, exists := sm.cacheService.Get(ctx, "health_check"); !exists {
			// Try to set and get a test value
			err := sm.cacheService.Set(ctx, "health_check", "ok", time.Minute)
			if err != nil {
				cacheHealth = "unhealthy"
			} else {
				_, exists := sm.cacheService.Get(ctx, "health_check")
				if !exists {
					cacheHealth = "unhealthy"
				}
			}
		}
		health["cache"] = cacheHealth
	} else {
		health["cache"] = "not_initialized"
	}

	// Check Job Queue
	if sm.jobQueue != nil {
		queueHealth := "healthy"
		stats, err := sm.jobQueue.GetStats(ctx)
		if err != nil {
			queueHealth = "unhealthy"
		} else {
			health["job_queue_stats"] = stats
		}
		health["job_queue"] = queueHealth
	} else {
		health["job_queue"] = "not_initialized"
	}

	// Check Rate Limit Service
	if sm.rateLimitService != nil {
		rateLimitHealth := "healthy"
		remaining, _ := sm.rateLimitService.GetRemaining(ctx, 1, "127.0.0.1")
		if remaining < 0 {
			rateLimitHealth = "unhealthy"
		}
		health["rate_limit"] = rateLimitHealth
	} else {
		health["rate_limit"] = "not_initialized"
	}

	health["services_initialized"] = sm.initialized

	return health
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
