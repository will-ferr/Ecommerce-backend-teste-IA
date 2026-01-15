package tests

import (
	"smart-choice/config"
	"smart-choice/database"
	"smart-choice/logger"
	"testing"

	"github.com/joho/godotenv"
)

// TestSetup initializes test environment
func TestSetup(t *testing.T) {
	// Load test environment variables
	if err := godotenv.Load("../.env.test"); err != nil {
		// Fallback to default .env if test env doesn't exist
		if err := godotenv.Load("../.env"); err != nil {
			t.Log("Warning: No .env file found, using environment variables")
		}
	}

	// Initialize logger (suppress output in tests)
	logger.InitLogger()

	// Load configuration
	config.LoadEnv()

	// Connect to test database
	database.ConnectDB()
}

// TestCleanup cleans up test environment
func TestCleanup(t *testing.T) {
	// Close database connection
	if database.DB != nil {
		if sqlDB, err := database.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
