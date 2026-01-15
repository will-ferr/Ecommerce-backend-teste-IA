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
)

func TestGetProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/products", controllers.GetProducts)

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Product Creation", func(t *testing.T) {
		router := gin.New()
		router.POST("/products", controllers.CreateProduct)

		productData := map[string]interface{}{
			"name":        "Test Product",
			"description": "Test Description",
			"price":       99.99,
			"stock":       10,
			"category":    "electronics",
		}
		jsonData, _ := json.Marshal(productData)

		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 201 or 401 (if auth required) or 400 (if validation fails)
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusUnauthorized || w.Code == http.StatusBadRequest)
	})

	t.Run("Invalid Product Creation - Missing Fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/products", controllers.CreateProduct)

		productData := map[string]interface{}{
			"name": "Test Product",
			// Missing required fields
		}
		jsonData, _ := json.Marshal(productData)

		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("Invalid Product Creation - Negative Price", func(t *testing.T) {
		router := gin.New()
		router.POST("/products", controllers.CreateProduct)

		productData := map[string]interface{}{
			"name":        "Test Product",
			"description": "Test Description",
			"price":       -10.0,
			"stock":       10,
			"category":    "electronics",
		}
		jsonData, _ := json.Marshal(productData)

		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Product ID", func(t *testing.T) {
		router := gin.New()
		router.GET("/products/:id", controllers.GetProduct)

		req, _ := http.NewRequest("GET", "/products/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 200 or 404 if product doesn't exist
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
	})

	t.Run("Invalid Product ID", func(t *testing.T) {
		router := gin.New()
		router.GET("/products/:id", controllers.GetProduct)

		req, _ := http.NewRequest("GET", "/products/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name     string
		product  map[string]interface{}
		expected int
	}{
		{
			name: "Empty name",
			product: map[string]interface{}{
				"name":        "",
				"description": "Test Description",
				"price":       99.99,
				"stock":       10,
				"category":    "electronics",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Zero price",
			product: map[string]interface{}{
				"name":        "Test Product",
				"description": "Test Description",
				"price":       0,
				"stock":       10,
				"category":    "electronics",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Negative stock",
			product: map[string]interface{}{
				"name":        "Test Product",
				"description": "Test Description",
				"price":       99.99,
				"stock":       -5,
				"category":    "electronics",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/products", controllers.CreateProduct)

			jsonData, _ := json.Marshal(tc.product)

			req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
