package repository

import (
	"strconv"

	"smart-choice/database"
	"smart-choice/models"
	"smart-choice/utils"
)

func GetProductByID(id uint) (models.Product, error) {
	var product models.Product
	err := database.DB.First(&product, id).Error
	return product, err
}

func GetProducts(pagination *utils.Pagination) ([]models.Product, error) {
	var products []models.Product
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := database.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuider.Model(&models.Product{}).Find(&products)
	return products, result.Error
}

func GetProductsWithFilters(pagination *utils.Pagination, name, minPrice, maxPrice, inStock string) ([]models.Product, error) {
	var products []models.Product
	query := database.DB.Model(&models.Product{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", min)
		}
	}

	if maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", max)
		}
	}

	if inStock == "true" {
		query = query.Where("stock > 0")
	}

	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Find(&products).Error
	return products, err
}
