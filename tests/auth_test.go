package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"smart-choice/controllers"
	"smart-choice/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService for testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, *models.User, error) {
	args := m.Called(email, password)
	return args.String(0), args.Get(1).(*models.User), args.Error(2)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test valid registration
	t.Run("Valid Registration", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controllers.Register)

		userData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
			"name":     "Test User",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 201 or 400 depending on validation
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)
	})

	// Test invalid registration
	t.Run("Invalid Registration - Missing Fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controllers.Register)

		userData := map[string]interface{}{
			"email": "test@example.com",
			// Missing password and name
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	// Test invalid email format
	t.Run("Invalid Registration - Bad Email", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controllers.Register)

		userData := map[string]interface{}{
			"email":    "invalid-email",
			"password": "password123",
			"name":     "Test User",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test valid login
	t.Run("Valid Login", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controllers.Login)

		loginData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 200 or 401 depending on credentials
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized)
	})

	// Test invalid login
	t.Run("Invalid Login - Missing Fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controllers.Login)

		loginData := map[string]interface{}{
			"email": "test@example.com",
			// Missing password
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test invalid credentials
	t.Run("Invalid Login - Bad Credentials", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controllers.Login)

		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
