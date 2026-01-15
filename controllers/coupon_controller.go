package controllers

import (
	"net/http"

	"smart-choice/services"

	"github.com/gin-gonic/gin"
)

type ValidateCouponInput struct {
	Code   string  `json:"code" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func ValidateCoupon(c *gin.Context) {
	var input ValidateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon, err := services.ValidateCoupon(input.Code, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupon": coupon})
}
