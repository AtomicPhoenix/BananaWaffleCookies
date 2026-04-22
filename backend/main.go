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
		r.Put("/api/profile/preferences", handlers.UpdateProfilePreferences)
		r.Get("/api/profile/preferences", handlers.GetProfilePreferences)
		r.Post("/api/profile/education", handlers.AddProfileEducation)
		r.Get("/api/profile/education", handlers.GetProfileEducation)
		r.Delete("/api/profile/education/{id}", handlers.DeleteProfileEducation)
		r.Put("/api/profile/education", handlers.UpdateProfileEducation)
		r.Put("/api/profile/education/reorder", handlers.ReorderProfileEducation)
		r.Post("/api/profile/experiences", handlers.AddProfileExperience)
		r.Get("/api/profile/experiences", handlers.GetProfileExperiences)
		r.Delete("/api/profile/experiences/{id}", handlers.DeleteProfileExperience)
		r.Put("/api/profile/experiences/reorder", handlers.ReorderProfileExperiences)
		r.Put("/api/profile/experiences", handlers.UpdateProfileExperience)
		r.Post("/api/profile/skills", handlers.AddProfileSkill)
		r.Put("/api/profile/skills", handlers.UpdateProfileSkill)
		r.Get("/api/profile/skills", handlers.GetProfileSkills)
		r.Delete("/api/profile/skills/{id}", handlers.DeleteProfileSkill)
		r.Put("/api/profile/skills/reorder", handlers.ReorderProfileSkill)
		r.Route("/api/jobs", func(r chi.Router) {
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

	portStr := fmt.Sprintf(":%d", settings.CLIArgs.GetPort())

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))
}
