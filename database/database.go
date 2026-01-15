package database

import (
	"fmt"
	"os"
	"smart-choice/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.Coupon{}, &models.ActivityLog{})
}

func ConnectDB() {
	var err error
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal().Msg("Failed to connect to database")
	}

	autoMigrate(DB)
}
