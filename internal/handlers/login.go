package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BBaCode/pocketwise-server/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleUserLogin(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var creds models.Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Email:", creds.Email)
	fmt.Println("Password:", creds.Password)

	if len(creds.Email) > 0 && len(creds.Password) > 0 {
		query := `SELECT COUNT(*) FROM public.users WHERE email=$1 AND password_hash=$2`

		var count int
		err := pool.QueryRow(context.Background(), query, creds.Email, creds.Password).Scan(&count)
		if err != nil {
			http.Error(w, "Failed to query user", http.StatusInternalServerError)
			return
		}

		// If the user exists (count > 0), return a success response
		if count > 0 {
			resp := models.Response{
				Status:  "Success",
				Message: "User with email " + creds.Email + " successfully logged in",
			}

			jsonData, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			// If no user is found, return a not found response
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		}
	} else {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
	}
}
