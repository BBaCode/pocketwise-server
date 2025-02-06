package app

import (
	"log"
	"os"

	"github.com/BBaCode/pocketwise-server/models"
)

func LoadConfig() models.DBConfig {
	// Get the Supabase connection string
	connString := os.Getenv("SUPABASE_DB_URL")
	if connString == "" {
		log.Fatalf("SUPABASE_DB_URL is not set in .env")
	}

	return models.DBConfig{
		ConnectionString: connString,
	}
}
