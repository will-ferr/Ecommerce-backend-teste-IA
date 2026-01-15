package repository

import (
	"smart-choice/database"
	"smart-choice/models"
	"time"
)

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return user, err
}

func CountNewUsers(start, end time.Time) (int64, error) {
	var count int64
	err := database.DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error
	return count, err
}
