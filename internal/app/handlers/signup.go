package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleUserSignUp(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := map[string]bool{
		"https://deploy-preview-13--pocketwise.netlify.app": true,
		"https://pocketwise.netlify.app":                    true,
	}

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON input
	var userDetails models.SignupRequest
	err = json.Unmarshal(body, &userDetails)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	query := `INSERT INTO public.users (id, email, first_name, last_name) VALUES ($1, $2, $3, $4)`
	_, err = pool.Exec(context.Background(), query, userDetails.Id, userDetails.Email, userDetails.FirstName, userDetails.LastName)
	if err != nil {
		fmt.Printf("Unexpected error: %s", err)
		http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
		return
	}

	// Success response
	resp := models.Response{
		Status:  "Success",
		Message: fmt.Sprintf("User with email %s successfully signed up", userDetails.Email),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
