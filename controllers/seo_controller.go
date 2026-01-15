package controllers

import (
	"net/http"
	"strconv"

	"smart-choice/services"

	"github.com/gin-gonic/gin"
)

func GetProductMetaTags(c *gin.Context) {
	productIDStr := c.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	metaTags, err := services.GetProductMetaTags(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, metaTags)
}
