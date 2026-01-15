package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	mu      sync.Mutex
	clients = make(map[string]*rate.Limiter)
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second
		}
		if !clients[ip].Allow() {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		mu.Unlock()
		c.Next()
	}
}
