package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bananawafflecookies.com/m/v2/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
		r.Use(jwtauth.Verifier(handlers.AuthToken))
		r.Use(jwtauth.Authenticator(handlers.AuthToken))

		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("Hello user %v", claims["id"])))
		})
	})

	router.Handle("/*", http.FileServer(http.Dir("./frontend/dist")))
	portStr := fmt.Sprintf(":%d", *config.port)

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))
}
