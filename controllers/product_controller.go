package controllers

import (
	"net/http"
	"strconv"

	"smart-choice/database"
	"smart-choice/models"
	"smart-choice/repository"
	"smart-choice/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       uint    `json:"stock" binding:"required"`
	StockLimit  uint    `json:"stock_limit"`
}

func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	nameErrors := utils.ValidateProductName(req.Name)
	priceErrors := utils.ValidatePrice(req.Price)
	stockErrors := utils.ValidateStock(req.Stock)

	allErrors := append(nameErrors, priceErrors...)
	allErrors = append(allErrors, stockErrors...)

	if len(allErrors) > 0 {
		utils.HandleValidationError(c, allErrors)
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		StockLimit:  req.StockLimit,
	}

	if req.StockLimit == 0 {
		product.StockLimit = 5
	}

	if err := database.DB.Create(&product).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func GetProducts(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)

	name := c.Query("name")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	inStock := c.Query("in_stock")

	products, err := repository.GetProductsWithFilters(&pagination, name, minPrice, maxPrice, inStock)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.StockLimit > 0 {
		product.StockLimit = req.StockLimit
	}

	if err := database.DB.Save(&product).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := database.DB.Delete(&models.Product{}, uint(id)).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
