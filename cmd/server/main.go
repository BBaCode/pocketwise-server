package main

import (
	"log"
	"net/http"

	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/internal/handlers"
	"github.com/BBaCode/pocketwise-server/lib/config"
)

func main() {
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()

	// Connect to the database
	pool, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	// Set up handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRequest(w, r, pool)
	})

	log.Println("Server starting on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("ListenAndServe failed: %v\n", err)
	}
}
