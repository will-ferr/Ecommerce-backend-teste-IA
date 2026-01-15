package services

import (
	"errors"
	"time"

	"smart-choice/models"
	"smart-choice/repository"
)

func ValidateCoupon(code string, amount float64) (*models.Coupon, error) {
	coupon, err := repository.GetCouponByCode(code)
	if err != nil {
		return nil, errors.New("invalid coupon code")
	}

	if coupon.ValidUntil.Before(time.Now()) {
		return nil, errors.New("coupon has expired")
	}

	if coupon.UsedCount >= coupon.MaxUses {
		return nil, errors.New("coupon has reached its usage limit")
	}

	if amount < coupon.MinAmount {
		return nil, errors.New("order amount does not meet the minimum requirement for this coupon")
	}

	return &coupon, nil
}
