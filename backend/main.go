package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/handlers"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func initEnvs() {
	godotenv.Load("./.env")

	// Confirm existence of necessary environmental variables
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("GEMINI_API_KEY env not present")
	}
}

func initDirs() {
	// Create data folder
	err := os.MkdirAll("data", 0750)
	if err != nil && err != fs.ErrExist {
		log.Fatalf("Failed to create data directory: %s\n", err)
	}
	err = os.MkdirAll("data/documents", 0750)
	if err != nil && err != fs.ErrExist {
		log.Fatalf("Failed to create data/documents directory: %s\n", err)
	}
	err = os.MkdirAll("data/logs", 0750)
	if err != nil && err != fs.ErrExist {
		log.Fatalf("Failed to create data/logs directory: %s\n", err)
	}
}

func init() {
	initEnvs()
	initDirs()

	settings.InitArgs()
	settings.InitLogs()
	handlers.InitAuth()

	// Initialize DB
	err := db.InitDB()
	if err != nil {
		log.Fatalf(`Failed to init database: %v`, err)
	}
}

func main() {
	router := chi.NewRouter()

	// API Routes
	router.Route("/api", func(r chi.Router) {
		// Public API Routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", handlers.RegistrationHandler)
			r.Post("/login", handlers.LoginHandler)
			r.Post("/logout", handlers.LogoutHandler)
		})

		// Protected API routes
		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware)

			// Auth check
			r.Get("/auth", handlers.GetAuth)

			// Profile API
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", handlers.GetProfile)
				r.Put("/", handlers.UpdateProfile)

				r.Route("/preferences", func(r chi.Router) {
					r.Get("/", handlers.GetProfilePreferences)
					r.Put("/", handlers.UpdateProfilePreferences)
				})

				r.Route("/education", func(r chi.Router) {
					r.Get("/", handlers.GetProfileEducation)
					r.Post("/", handlers.AddProfileEducation)
					r.Put("/", handlers.UpdateProfileEducation)
					r.Put("/reorder", handlers.ReorderProfileEducation)
					r.Delete("/{id}", handlers.DeleteProfileEducation)
				})

				r.Route("/experiences", func(r chi.Router) {
					r.Get("/", handlers.GetProfileExperiences)
					r.Post("/", handlers.AddProfileExperience)
					r.Put("/", handlers.UpdateProfileExperience)
					r.Put("/reorder", handlers.ReorderProfileExperiences)
					r.Delete("/{id}", handlers.DeleteProfileExperience)
				})

				r.Route("/skills", func(r chi.Router) {
					r.Get("/", handlers.GetProfileSkills)
					r.Post("/", handlers.AddProfileSkill)
					r.Put("/", handlers.UpdateProfileSkill)
					r.Put("/reorder", handlers.ReorderProfileSkill)
					r.Delete("/{id}", handlers.DeleteProfileSkill)
				})
			})

			// Jobs API
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", handlers.GetJobs)
				r.Post("/", handlers.CreateJob)
				r.Put("/", handlers.UpdateJob)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetJob)
					r.Delete("/", handlers.DeleteJob)

					r.Post("/archive", handlers.ArchiveJob)
					r.Post("/unarchive", handlers.UnarchiveJob)

					r.Post("/resume", handlers.GetResumeDraft)
					r.Post("/cover-letter", handlers.GetCoverLetterDraft)

					r.Get("/activities", handlers.GetJobActivities)

					r.Route("/interviews", func(r chi.Router) {
						r.Get("/", handlers.GetInterviews)
						r.Post("/", handlers.CreateInterview)
						r.Delete("/{interview_id}", handlers.DeleteInterview)
					})

					r.Route("/followups", func(r chi.Router) {
						r.Get("/", handlers.GetFollowUps)
						r.Post("/", handlers.CreateFollowUp)
						r.Put("/{followup_id}", handlers.UpdateFollowUp)
						r.Delete("/{followup_id}", handlers.DeleteFollowUp)
					})

					r.Route("/documents", func(r chi.Router) {
						r.Get("/", handlers.GetJobDocuments)
						r.Post("/", handlers.LinkDocumentToJob)
						r.Delete("/{document_id}", handlers.UnlinkDocumentFromJob)
						r.Post("/ai-save", handlers.SaveAIDocumentToJob)
					})
				})
			})

			// Documents API
			r.Route("/documents", func(r chi.Router) {
				r.Get("/{id}", handlers.GetDocument)
				r.Get("/{id}/info", handlers.GetDocumentInfo)
				r.Post("/", handlers.UploadDocument)
				r.Put("/{id}", handlers.UpdateDocument)
				r.Delete("/{id}", handlers.DeleteDocument)
			})

			// Settings API
			r.Route("/settings", func(r chi.Router) {
				r.Get("/", handlers.GetSettings)
				r.Put("/", handlers.UpdateSettings)
			})
		})
	})

	// Frontend / Static Routes
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	})

	// Serve frontend for Vue routes
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	})

	// Serve static images
	router.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/dist/images"))))
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/dist/assets"))))

	portStr := fmt.Sprintf(":%d", settings.CLIArgs.GetPort())

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))
}
