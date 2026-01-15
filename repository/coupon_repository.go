package repository

import (
	"smart-choice/database"
	"smart-choice/models"
)

func GetCouponByCode(code string) (models.Coupon, error) {
	var coupon models.Coupon
	err := database.DB.Where("code = ?", code).First(&coupon).Error
	return coupon, err
}

func UpdateCoupon(coupon *models.Coupon) error {
	return database.DB.Save(coupon).Error
}
