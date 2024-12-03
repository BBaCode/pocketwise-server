package main

import (
	"log"
	"net/http"

	config "github.com/BBaCode/pocketwise-server/internal/app"
	"github.com/BBaCode/pocketwise-server/internal/app/handlers"
	"github.com/BBaCode/pocketwise-server/internal/app/middleware"
	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	pool, err := db.Connect(db.DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUserSignUp(w, r, pool)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUserLogin(w, r, pool)
	}).Methods("POST", "OPTIONS")

	// Protected routes (With JWT validation)
	r.Handle("/accounts", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAccounts(w, r, pool)
	}))).Methods("GET", "OPTIONS")

	r.Handle("/new-accounts", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAddAccounts(w, r, pool)
	}))).Methods("GET")

	r.Handle("/transactions", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetTransactions(w, r, pool)
	}))).Methods("POST", "OPTIONS")

	log.Println("Server starting on :80")
	if err := http.ListenAndServe(":80", r); err != nil {
		log.Fatalf("ListenAndServe failed: %v\n", err)
	}

}
