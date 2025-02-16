package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BBaCode/pocketwise-server/internal/app/constants"
	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleGetBudget(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	var budgetRequest models.BudgetRequest
	err = json.Unmarshal(body, &budgetRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var budgetResponse models.StoredBudget
	budgetResponse, err = db.FetchExistingBudget(budgetRequest, pool)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected error fetching accounts: %v", err), http.StatusInternalServerError)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(budgetResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func HandleGetAllBudgets(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var storedBudgets []models.StoredBudget
	storedBudgets, err := db.FetchAllExistingBudgets(userID, pool)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected error fetching accounts: %v", err), http.StatusInternalServerError)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(storedBudgets); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func HandleAddNewBudget(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the raw body to debug
	fmt.Println("Raw request body:", string(body))

	// Parse the JSON input
	var budgetRequest models.BudgetRequest
	err = json.Unmarshal(body, &budgetRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	budgetRequest.UserId = userID

	var budgetResponse models.MessageResponse
	err = db.InsertNewBudget(budgetRequest, pool)
	if err != nil {
		fmt.Print(err)
		if err.Error() == `ERROR: duplicate key value violates unique constraint "unique_month_year" (SQLSTATE 23505)` {
			budgetResponse.Message = "Budget already exists for that month/year."
		} else {
			budgetResponse.Message = "Budget could not be created, please try again later."
		}
	} else {
		budgetResponse.Message = "Budget created successfully"
	}
	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(budgetResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func HandleDeleteBudget(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	budgetId := vars["budgetId"]

	if budgetId == "" {
		http.Error(w, "Budget ID is required", http.StatusBadRequest)
		return
	}

	// Parse the JSON input

	var budgetResponse models.MessageResponse
	if db.DeleteBudget(budgetId, pool) == nil {
		budgetResponse.Message = "Budget deleted successfully"

	} else {
		budgetResponse.Message = "Budget could not be deleted, please try again later."
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(budgetResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func HandleUpdateBudget(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	vars := mux.Vars(r)
	budgetId := vars["budgetId"]

	if budgetId == "" {
		http.Error(w, "Budget ID is required", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the raw body to debug
	fmt.Println("Raw request body:", string(body))

	// Parse the JSON input
	var updateBudgetRequest models.UpdateBudgetRequest
	err = json.Unmarshal(body, &updateBudgetRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var budgetResponse models.MessageResponse
	if db.UpdateExistingBudget(budgetId, updateBudgetRequest, pool) == nil {
		budgetResponse.Message = "Budget updated successfully"
	} else {
		budgetResponse.Message = "Budget could not be updated, please try again later."
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(budgetResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}
