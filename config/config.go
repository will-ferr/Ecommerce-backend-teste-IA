package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
}
