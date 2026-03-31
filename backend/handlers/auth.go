package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bananawafflecookies.com/m/v2/db"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET_KEY string
var AuthToken *jwtauth.JWTAuth

func init() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	AuthToken = jwtauth.New("HS256", []byte(JWT_SECRET_KEY), nil)

}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handles user registration
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	password_bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		http.Error(w, "Failed to generate password hash", http.StatusInternalServerError)
		return
	}
	password_hash := string(password_bytes)

	user := db.User{
		Email:         req.Email,
		Password_hash: password_hash,
	}
	_, err = db.RegisterUser(user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	user, err := db.GetUser(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		fmt.Printf("Failed to encrypt password: %v", err)
		return
	}

	// Create JWT token
	claims := map[string]interface{}{
		"id":    fmt.Sprintf("%v", user.Id),
		"email": user.Email,
	}

	_, tokenString, err := AuthToken.Encode(claims)
	if err != nil {
		fmt.Println("JWT encode error:", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set JWT as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Respond to client with success message
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"Login successful."}`)
}

// Handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Overwrite the auth_token cookie with empty value and expired date
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
}
