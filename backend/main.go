package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bananawafflecookies.com/m/v2/handlers"
	"github.com/go-chi/chi/v5"
)

// CLI Arguments
type Config struct {
	dev  *bool
	port *int
}

var config Config

func setup() {
	// Parse CLI Arguments
	config.dev = flag.Bool("dev", false, "run in development mode")
	config.port = flag.Int("p", 8080, "port to run server on")
	flag.Parse()
}

func main() {
	setup()

	router := chi.NewRouter()

	// Public routes
	router.Post("/signup", handlers.RegistrationHandler)
	router.Post("/login", handlers.LoginHandler)
	router.Post("/logout", handlers.LogoutHandler)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(AuthRedirect)

		r.Get("/api/auth/check", func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			w.Header().Set("Content-Type", "application/json")
			if err != nil || cookie.Value == "" {
				w.Write([]byte(`{"authenticated": false}`))
				return
			}

			w.Write([]byte(`{"authenticated": true}`))
		})

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {})
		r.Put("/api/profile", handlers.UpdateProfile)
		r.Get("/api/profile", handlers.GetProfile)

		r.Get("/library", func(w http.ResponseWriter, r *http.Request) {})
		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {})
		r.Get("/settings", func(w http.ResponseWriter, r *http.Request) {})
		r.Put("/api/settings", func(w http.ResponseWriter, r *http.Request) {})
		r.Get("/api/settings", handlers.GetSettings)
		r.Get("/create-job", func(w http.ResponseWriter, r *http.Request) {})
	})

	// Serve frontend for Vue routes
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	})
	// Serve static assets
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/dist/assets"))))

	portStr := fmt.Sprintf(":%d", *config.port)

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))
}

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			// If not authenticated, redirect to /login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = handlers.AuthToken.Decode(cookie.Value)
		if err != nil {
			// If not authenticated, redirect to /login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
