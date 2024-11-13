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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, OPTIONS")
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

	// we're setting this before doing the insert so would need to be changed in case of a fail
	resp := models.Response{
		Status:  "Success",
		Message: "User with email " + creds.Email + " successfully signed up",
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	if len(creds.Email) > 0 && len(creds.Password) > 0 {
		query := `INSERT INTO public.users (email, password_hash) VALUES ($1, $2)`
		_, err = pool.Exec(context.Background(), query, creds.Email, creds.Password)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}
	}

	w.Write(jsonData)
}
