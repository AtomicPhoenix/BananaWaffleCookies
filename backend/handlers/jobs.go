package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bananawafflecookies.com/m/v2/db"
)

// Handler for /api/jobs (POST)
func CreateJob(w http.ResponseWriter, r *http.Request) {
	var job db.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to post job; Failed to grab decode request: %v\n", err)
		return
	}

	var tokenInfo Claim
	err, tokenInfo = GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to post job; Failed to grab auth token information: %v\n", err)
		return
	}

	job.UserID = tokenInfo.Uid
	_, err = db.CreateJob(job)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to post job: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully uploaded."}`)
}

// Handler for /api/jobs (GET)
func GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := db.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get jobs: %v\n", err)
		return
	}
	if jobs == nil {
		jobs = []db.Job{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
