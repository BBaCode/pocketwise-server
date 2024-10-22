package main

import (
	"log"
	"net/http"

	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/internal/handlers"
	"github.com/BBaCode/pocketwise-server/lib/config"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()

	// Connect to the database
	pool, err := db.Connect(db.DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	r := mux.NewRouter()

	http.Handle("/", r)

	// Set up handlers
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUserSignUp(w, r, pool)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUserLogin(w, r, pool)
	}).Methods("POST")

	log.Println("Server starting on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("ListenAndServe failed: %v\n", err)
	}
}
