package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"bananawafflecookies.com/m/v2/db"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET_KEY string
var AuthToken *jwtauth.JWTAuth

type Claim struct {
	Uid   int
	Email string
}

func InitAuth() {
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	if JWT_SECRET_KEY == "" {
		log.Fatal("Environmental Variable JWT_SECRET_KEY is not set")
	}
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

// Grabs information from token claims
func GrabToken(r *http.Request) (error, Claim) {
	cookie, err := r.Cookie("auth_token")

	// Grab cookie
	if err != nil || cookie.Value == "" {
		return err, Claim{}
	}

	// Decode auth token
	reqToken, err := AuthToken.Decode(cookie.Value)
	if err != nil {
		fmt.Println("Failed to decode reqToken:", err)
		return err, Claim{}
	}

	// Auth Token fields to grab
	var uid_str, email string

	// Decode auth token
	err = reqToken.Get("id", &uid_str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode user id from auth token: %v\n", err)
		return err, Claim{}
	}

	uid, err := strconv.Atoi(uid_str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return err, Claim{}
	}

	// Grab email
	err = reqToken.Get("email", &email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode email from auth token: %v\n", err)
		return err, Claim{}
	}

	return nil, Claim{Uid: uid, Email: email}
}

func GetAuth(w http.ResponseWriter, r *http.Request) {
	err, _ := GrabToken(r)

	w.Header().Set("Content-Type", "application/json")

	// GrabToken will error if it fails to find a valid authorization token
	if err != nil {
		w.Write([]byte(`{"authenticated": false}`))
		return
	}
	w.Write([]byte(`{"authenticated": true}`))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, claim := GrabToken(r)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}

		// Update request context to store user info
		ctx := context.WithValue(r.Context(), "user", claim)

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
