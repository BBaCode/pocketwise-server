package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/golang-jwt/jwt/v5" // Import JWT library
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

	// Read the request body
	body, err := io.ReadAll(r.Body)
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

		// Generate a JWT token upon successful login
		// token, err := generateJWT(userDetails.Email)
		// if err != nil {
		// 	fmt.Print(err)
		// 	http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		// 	return
		// }

		// Send the token in the response
		resp := models.Response{
			Status:  "Success",
			Message: "User successfully logged in",
			// Token:   token, // Include the token in your response model
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

}

// Generate a JWT token
func generateJWT(email string) (string, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the JWT secret from the environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in the environment")
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token valid for 1 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
