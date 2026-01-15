package repository

import (
	"time"

	"smart-choice/database"
	"smart-choice/models"
)

func GetTotalSales(start, end time.Time) (float64, error) {
	var total float64
	err := database.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", start, end).Select("sum(total)").Row().Scan(&total)
	return total, err
}

func CreateOrder(order *models.Order) error {
	return database.DB.Create(order).Error
}

func UpdateOrderStatus(orderID uint, status string) error {
	return database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func GetOrderStatusCounts() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}
	counts := make(map[string]int64)

	err := database.DB.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}
