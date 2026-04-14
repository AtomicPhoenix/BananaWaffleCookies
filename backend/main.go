package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

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

	// Create data folder
	err = os.MkdirAll("data", 0750)
	if err != nil && err != fs.ErrExist {
		log.Fatalf("Failed to create data directory: %s\n", err)
	}
	err = os.MkdirAll("data/documents", 0750)
	if err != nil && err != fs.ErrExist {
		log.Fatalf("Failed to create data/documents directory: %s\n", err)
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
		r.Post("/api/profile/education", handlers.AddProfileEducation)
		r.Get("/api/profile/education", handlers.GetProfileEducation)
		r.Delete("/api/profile/education/{id}", handlers.DeleteProfileEducation)
		r.Put("/api/profile/education/reorder", handlers.ReorderProfileEducation)
		r.Post("/api/profile/experiences", handlers.AddProfileExperience)
		r.Get("/api/profile/experiences", handlers.GetProfileExperiences)
		r.Delete("/api/profile/experiences/{id}", handlers.DeleteProfileExperience)
		r.Put("/api/profile/experiences/reorder", handlers.ReorderProfileExperiences)
		r.Post("/api/profile/skills", handlers.AddProfileSkill)
		r.Get("/api/profile/skills", handlers.GetProfileSkills)
		r.Delete("/api/profile/skills/{id}", handlers.DeleteProfileSkill)
		r.Put("/api/profile/skills/reorder", handlers.ReorderProfileSkill)
		r.Get("/api/jobs", handlers.GetJobs)
		r.Get("/api/jobs/{id}", handlers.GetJob)
		r.Post("/api/jobs", handlers.CreateJob)
		r.Put("/api/jobs", handlers.UpdateJob)
		r.Put("/api/settings", handlers.UpdateSettings)
		r.Get("/api/settings", handlers.GetSettings)
		r.Get("/api/documents/{id}", handlers.GetDocument)
		r.Get("/api/documents/{id}/info", handlers.GetDocumentInfo)
		r.Post("/api/documents", handlers.UploadDocument)
		r.Delete("/api/documents/{id}", handlers.DeleteDocument)
		r.Put("/api/documents/{id}", handlers.UpdateDocument)
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
