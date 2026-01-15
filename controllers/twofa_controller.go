package controllers

import (
	"net/http"

	"smart-choice/models"
	"smart-choice/services"

	"github.com/gin-gonic/gin"
)

func Generate2FA(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user := userCtx.(*models.User)
	qrCode, secret, err := services.Generate2FA(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": qrCode, "secret": secret})
}

type Validate2FAInput struct {
	Code string `json:"code" binding:"required"`
}

func Validate2FA(c *gin.Context) {
	var input Validate2FAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user := userCtx.(*models.User)
	if !services.Validate2FA(user, input.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	token, err := services.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
