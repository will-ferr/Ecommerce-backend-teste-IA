package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"smart-choice/services"

	"github.com/gin-gonic/gin"
)

type EnhancedRateLimiter struct {
	rateLimitService services.RateLimitService
}

func NewEnhancedRateLimiter(rateLimitService services.RateLimitService) *EnhancedRateLimiter {
	return &EnhancedRateLimiter{
		rateLimitService: rateLimitService,
	}
}

func (e *EnhancedRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			userID = uint(0) // Anonymous user
		}

		// Get client IP
		clientIP := c.ClientIP()

		// Check rate limit
		ctx := context.Background()
		allowed := e.rateLimitService.Allow(ctx, userID.(uint), clientIP)

		if !allowed {
			// Get remaining requests for headers
			remaining, resetTime := e.rateLimitService.GetRemaining(ctx, userID.(uint), clientIP)

			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(resetTime).Unix(), 10))
			c.Header("Retry-After", strconv.Itoa(int(resetTime.Seconds())))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		remaining, resetTime := e.rateLimitService.GetRemaining(ctx, userID.(uint), clientIP)
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(resetTime).Unix(), 10))

		c.Next()
	}
}
