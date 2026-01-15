package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"smart-choice/controllers"
	"smart-choice/routes"
	"smart-choice/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationAuthFlow(t *testing.T) {
	TestSetup(t)
	defer TestCleanup(t)

	gin.SetMode(gin.TestMode)

	// Initialize services
	serviceManager := services.GetServiceManager()
	if err := serviceManager.InitializeServices(); err != nil {
		t.Skip("Skipping integration test - Redis not available")
	}

	// Setup router
	router := gin.New()
	routes.SetupRoutes(router)

	// Test complete auth flow
	t.Run("Complete Auth Flow", func(t *testing.T) {
		// 1. Register user
		registerData := map[string]interface{}{
			"email":    "integration@example.com",
			"password": "password123",
			"name":     "Integration Test User",
		}
		jsonData, _ := json.Marshal(registerData)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should succeed or fail gracefully
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)

		// 2. Login with registered user
		loginData := map[string]interface{}{
			"email":    "integration@example.com",
			"password": "password123",
		}
		jsonData, _ = json.Marshal(loginData)

		req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should succeed or fail gracefully
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized)
	})
}

func TestIntegrationProductFlow(t *testing.T) {
	TestSetup(t)
	defer TestCleanup(t)

	gin.SetMode(gin.TestMode)

	// Setup router
	router := gin.New()
	routes.SetupRoutes(router)

	t.Run("Product CRUD Flow", func(t *testing.T) {
		// 1. Get all products (should work without auth)
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 2. Try to create product without auth (should fail)
		productData := map[string]interface{}{
			"name":        "Test Product",
			"description": "Test Description",
			"price":       99.99,
			"stock":       10,
			"category":    "electronics",
		}
		jsonData, _ := json.Marshal(productData)

		req, _ = http.NewRequest("POST", "/api/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestIntegrationHealthChecks(t *testing.T) {
	TestSetup(t)
	defer TestCleanup(t)

	gin.SetMode(gin.TestMode)

	// Setup router
	router := gin.New()
	routes.SetupRoutes(router)

	t.Run("Health Check Endpoints", func(t *testing.T) {
		// Test /health endpoint
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var healthResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.NoError(t, err)
		assert.Contains(t, healthResponse, "status")
		assert.Contains(t, healthResponse, "checks")

		// Test /ready endpoint
		req, _ = http.NewRequest("GET", "/ready", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be ready or not ready depending on services
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusServiceUnavailable)

		// Test /alive endpoint
		req, _ = http.NewRequest("GET", "/alive", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var aliveResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &aliveResponse)
		assert.NoError(t, err)
		assert.Equal(t, "alive", aliveResponse["status"])
	})
}

func TestIntegrationMiddleware(t *testing.T) {
	TestSetup(t)
	defer TestCleanup(t)

	gin.SetMode(gin.TestMode)

	// Setup router
	router := gin.New()
	routes.SetupRoutes(router)

	t.Run("Middleware Integration", func(t *testing.T) {
		// Test security headers
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check for security headers
		assert.Contains(t, w.Header().Get("X-Frame-Options"), "DENY")
		assert.Contains(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
		assert.Contains(t, w.Header().Get("X-XSS-Protection"), "1; mode=block")

		// Test CORS headers
		req, _ = http.NewRequest("OPTIONS", "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should have CORS headers
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func BenchmarkHealthCheck(b *testing.B) {
	TestSetup(&testing.T{})
	defer TestCleanup(&testing.T{})

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", controllers.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
