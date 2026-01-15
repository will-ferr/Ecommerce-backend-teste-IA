package controllers

import (
	"context"
	"runtime"
	"smart-choice/database"
	"smart-choice/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]interface{} `json:"checks"`
	System    SystemInfo             `json:"system"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS           string     `json:"os"`
	Architecture string     `json:"architecture"`
	NumCPU       int        `json:"num_cpu"`
	NumGoroutine int        `json:"num_goroutine"`
	MemoryUsage  MemoryInfo `json:"memory_usage"`
}

// MemoryInfo represents memory usage information
type MemoryInfo struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
}

var startTime = time.Now()

// HealthCheck performs comprehensive health check
func HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   getAppVersion(),
		Uptime:    time.Since(startTime).String(),
		Checks:    make(map[string]interface{}),
		System:    getSystemInfo(),
	}

	// Database health check
	health.Checks["database"] = checkDatabaseHealth(ctx)

	// Redis services health check
	serviceManager := services.GetServiceManager()
	health.Checks["services"] = serviceManager.HealthCheck(ctx)

	// Overall status determination
	if !isHealthy(health.Checks) {
		health.Status = "unhealthy"
		c.JSON(503, health)
		return
	}

	c.JSON(200, health)
}

// ReadinessCheck checks if the application is ready to serve traffic
func ReadinessCheck(c *gin.Context) {
	ctx := c.Request.Context()

	ready := true
	checks := make(map[string]interface{})

	// Check database
	dbHealthy := checkDatabaseHealth(ctx)
	checks["database"] = dbHealthy
	if dbHealthy != "healthy" {
		ready = false
	}

	// Check Redis services
	serviceManager := services.GetServiceManager()
	servicesHealth := serviceManager.HealthCheck(ctx)
	checks["services"] = servicesHealth

	// Check if all services are healthy
	if servicesHealthMap, ok := interface{}(servicesHealth).(map[string]interface{}); ok {
		for service, status := range servicesHealthMap {
			if service != "job_queue_stats" && status != "healthy" {
				ready = false
				break
			}
		}
	} else {
		ready = false
	}

	status := "ready"
	if !ready {
		status = "not_ready"
		c.JSON(503, gin.H{
			"status":    status,
			"timestamp": time.Now(),
			"checks":    checks,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    status,
		"timestamp": time.Now(),
		"checks":    checks,
	})
}

// LivenessCheck checks if the application is alive
func LivenessCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime).String(),
	})
}

func checkDatabaseHealth(ctx context.Context) string {
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database instance")
		return "unhealthy"
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Error().Err(err).Msg("Database ping failed")
		return "unhealthy"
	}

	return "healthy"
}

func getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryUsage: MemoryInfo{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
		},
	}
}

func getAppVersion() string {
	// This could be set during build time with ldflags
	// Example: go build -ldflags "-X main.appVersion=1.0.0"
	return "1.0.0"
}

func isHealthy(checks map[string]interface{}) bool {
	for name, check := range checks {
		switch name {
		case "database":
			if status, ok := check.(string); ok && status != "healthy" {
				return false
			}
		case "services":
			if servicesHealth, ok := check.(map[string]interface{}); ok {
				for service, status := range servicesHealth {
					if service != "job_queue_stats" && status != "healthy" {
						return false
					}
				}
			} else {
				return false
			}
		}
	}
	return true
}
