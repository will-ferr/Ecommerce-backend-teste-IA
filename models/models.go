package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"-" gorm:"not null"`
	IsAdmin     bool   `json:"is_admin" gorm:"default:false"`
	TwoFA       bool   `json:"two_fa" gorm:"default:false"`
	TwoFASecret string `json:"-"`
}

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	StockLimit  uint    `json:"stock_limit" gorm:"default:5"`
}

func (p *Product) AfterUpdate(tx *gorm.DB) (err error) {
	if p.Stock < p.StockLimit {
		log := ActivityLog{
			Action:    fmt.Sprintf("Stock for product %s is low (%d remaining)", p.Name, p.Stock),
			Timestamp: time.Now(),
		}
		tx.Create(&log)
	}
	return
}

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
	Total      float64     `json:"total"`
	Status     string      `json:"status" gorm:"default:'pending'"`
	CouponID   *uint       `json:"coupon_id"`
	Coupon     *Coupon     `json:"coupon"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

type Coupon struct {
	gorm.Model
	Code       string    `json:"code" gorm:"unique"`
	Discount   float64   `json:"discount"`
	ValidUntil time.Time `json:"valid_until"`
	MaxUses    uint      `json:"max_uses"`
	UsedCount  uint      `json:"used_count" gorm:"default:0"`
	MinAmount  float64   `json:"min_amount"`
}

type ActivityLog struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}
