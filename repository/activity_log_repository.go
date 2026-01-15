package repository

import (
	"smart-choice/database"
	"smart-choice/models"
)

func CreateActivityLog(log *models.ActivityLog) error {
	return database.DB.Create(log).Error
}
