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
		fmt.Fprintf(os.Stderr, "Unable to load env file: %v\n", err)
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

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		http.Error(w, "Failed to generate password hash", http.StatusInternalServerError)
		return
	}
	password_hash := string(bytes)

	user := db.User{
		Email:         req.Email,
		Password_hash: password_hash,
	}
	_, err = db.RegisterUser(user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully created user")
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create & return user authorization token based on user data
	inputData := map[string]interface{}{"id": user.Id, "email": user.Email, "exp": time.Now().Add(24 * time.Hour).Unix()}
	_, tokenString, err := AuthToken.Encode(inputData)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, fmt.Sprintf(`{"token":"%s"}`, tokenString))
}

// Handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
