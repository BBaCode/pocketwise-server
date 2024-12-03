package app

import (
	"log"
	"os"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
)

func LoadConfig() models.DBConfig {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the Supabase connection string
	connString := os.Getenv("SUPABASE_DB_URL")
	if connString == "" {
		log.Fatalf("SUPABASE_DB_URL is not set in .env")
	}

	return models.DBConfig{
		ConnectionString: connString,
	}
}
