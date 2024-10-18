package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {

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

	var creds Credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Email:", creds.Email)
	fmt.Println("Password:", creds.Password)

	resp := Response{
		Status:  "success",
		Message: "You found me!",
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	if len(creds.Email) > 0 {
		query := `INSERT INTO public.users (email, password_hash) VALUES ($1, $2)`
		_, err = pool.Exec(context.Background(), query, creds.Email, creds.Password)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}
	}

	w.Write(jsonData)
}
