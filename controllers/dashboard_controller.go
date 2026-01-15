package controllers

import (
	"net/http"

	"smart-choice/services"

	"github.com/gin-gonic/gin"
)

func GetDashboardMetrics(c *gin.Context) {
	metrics, err := services.GetDashboardMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
