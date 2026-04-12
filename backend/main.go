package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// CLI Arguments
type Config struct {
	dev  *bool
	port *int
}

var config Config

func init() {
	// Parse CLI Arguments
	config.dev = flag.Bool("dev", false, "run in development mode")
	config.port = flag.Int("p", 8080, "port to run server on")
	flag.Parse()

	godotenv.Load("./.env")

	// Initiailize AuthToken
	handlers.InitAuth()

	// Initialize DB
	err := db.InitDB()
	if err != nil {
		log.Fatalf(`Failed to init database: %v`, err)
	}
}

func main() {
	router := chi.NewRouter()

	// Public API Routes
	router.Post("/api/signup", handlers.RegistrationHandler)
	router.Post("/api/login", handlers.LoginHandler)
	router.Post("/api/logout", handlers.LogoutHandler)

	// Protected API routes
	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)
		r.Get("/api/auth", handlers.GetAuth)
		r.Put("/api/profile", handlers.UpdateProfile)
		r.Get("/api/profile", handlers.GetProfile)
		r.Get("/api/jobs", handlers.GetJobs)
		r.Get("/api/jobs/{id}", handlers.GetJob)
		r.Post("/api/jobs", handlers.CreateJob)
		r.Put("/api/jobs", handlers.UpdateJob)
		r.Put("/api/settings", handlers.UpdateSettings)
		r.Get("/api/settings", handlers.GetSettings)
	})

	// Serve frontend for Vue routes
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	})

	// Serve static images
	router.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/dist/images"))))
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/dist/assets"))))

	portStr := fmt.Sprintf(":%d", *config.port)

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))
}
