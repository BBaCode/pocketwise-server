package main

import (
	"log"
	"net/http"
	"os"

	config "github.com/BBaCode/pocketwise-server/internal/app"
	"github.com/BBaCode/pocketwise-server/internal/app/handlers"
	"github.com/BBaCode/pocketwise-server/internal/app/middleware"
	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		// Try to load .env file only in local development
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Println("No .env file found, relying on system environment variables.")
		} else {
			log.Println(".env file loaded successfully.")
		}
	}

	cfg := config.LoadConfig()

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
	// Gets existing accounts from the database (no transactions)
	r.Handle("/accounts", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAccounts(w, r, pool)
	}))).Methods("GET", "OPTIONS")

	// this works for any NEW account but not for updating the same accounts.
	// currently this is done from an api, we dont have UI for this
	r.Handle("/new-accounts", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAddAccounts(w, r, pool)
	}))).Methods("GET")

	// Updates all accounts AND transactions: this would be the function to run every morning/ twice a day
	r.Handle("/account-data", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetUpdatedAccountData(w, r, pool)
	}))).Methods("GET", "OPTIONS")

	r.Handle("/all-transactions", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAllTransactions(w, r, pool)
	}))).Methods("GET", "POST", "OPTIONS")

	// this currently lets us load more data from simplefin into the transactions table by passing an account
	// not used at all as we moved away from account level updates
	r.Handle("/transactions", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetTransactions(w, r, pool)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/update-transactions", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUpdateTransactions(w, r, pool)
	}))).Methods("PUT", "OPTIONS")

	r.Handle("/budget", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetBudget(w, r, pool)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/all-budgets", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetAllBudgets(w, r, pool)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/new-budget", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAddNewBudget(w, r, pool)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/delete-budget/{budgetId}", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDeleteBudget(w, r, pool)
	}))).Methods("DELETE", "OPTIONS")

	r.Handle("/update-budget/{budgetId}", middleware.ValidateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUpdateBudget(w, r, pool)
	}))).Methods("PUT", "OPTIONS")

	log.Println("Server starting on :80")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("ListenAndServe failed: %v\n", err)
	}
}
