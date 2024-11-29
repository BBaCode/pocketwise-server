package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware to validate JWT
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := authHeader[len("Bearer "):]
		// Parse and validate the token
		secret := os.Getenv("SUPABASE_JWT_SECRET") // Load your Supabase JWT secret
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("Unexpected Error:%v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims (e.g., sub = user ID)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			r.Header.Set("X-User-ID", claims["sub"].(string))
		}

		next.ServeHTTP(w, r)
	})
}
