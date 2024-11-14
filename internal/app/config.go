package app

import (
	"log"
	"os"
	"strconv"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
)

func LoadConfig() models.DBConfig {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse the DB_PORT as an integer
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value in .env: %v", err)
	}

	// Return the Config struct populated with values from the environment
	return models.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}
