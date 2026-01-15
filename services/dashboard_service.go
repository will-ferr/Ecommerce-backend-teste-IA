package services

import (
	"time"

	"smart-choice/repository"
)

type DashboardMetrics struct {
	DailySales   float64          `json:"daily_sales"`
	MonthlySales float64          `json:"monthly_sales"`
	NewUsers     int64            `json:"new_users"`
	OrderStatus  map[string]int64 `json:"order_status"`
}

func GetDashboardMetrics() (*DashboardMetrics, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	dailySales, err := repository.GetTotalSales(startOfDay, now)
	if err != nil {
		return nil, err
	}

	monthlySales, err := repository.GetTotalSales(startOfMonth, now)
	if err != nil {
		return nil, err
	}

	newUsers, err := repository.CountNewUsers(startOfMonth, now)
	if err != nil {
		return nil, err
	}

	orderStatus, err := repository.GetOrderStatusCounts()
	if err != nil {
		return nil, err
	}

	metrics := &DashboardMetrics{
		DailySales:   dailySales,
		MonthlySales: monthlySales,
		NewUsers:     newUsers,
		OrderStatus:  orderStatus,
	}

	return metrics, nil
}
