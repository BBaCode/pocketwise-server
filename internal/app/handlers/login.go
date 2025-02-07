package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BBaCode/pocketwise-server/models" // Import JWT library
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleUserLogin(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON input
	var userDetails models.LoginRequest
	err = json.Unmarshal(body, &userDetails)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(userDetails.Id) == 0 {
		http.Error(w, "Cant find user", http.StatusBadRequest)
		return
	}

	// Fetch the stored hash for the provided email
	query := `SELECT email, first_name, last_name FROM public.users WHERE id=$1`
	var email, firstName, lastName string

	err = pool.QueryRow(context.Background(), query, userDetails.Id).Scan(&email, &firstName, &lastName)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	// Create a response object
	resp := models.Response{
		Status:  "Success",
		Message: "User successfully logged in",
		Data: models.LoginResponse{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		},
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
