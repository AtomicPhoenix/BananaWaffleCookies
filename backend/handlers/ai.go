package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"bananawafflecookies.com/m/v2/ai"
	"bananawafflecookies.com/m/v2/db"
	"github.com/go-chi/chi/v5"
)

// Handler for /api/jobs/{id}/resume (GET)
func GetResumeDraft(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to generate resume; Failed to grab auth token information: %v\n", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	job, err := db.GetJob(job_id, tokenInfo.Uid)

	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to generate resume: %v\n", err)
		return
	}

	profile, err := db.GetProfile(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Internal server error when generating resume", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return
	}

	response, err := ai.GenerateResumeDraft(job, profile)
	if err != nil {
		http.Error(w, "Internal server error when generating resume", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return
	}

	respJson := struct {
		Success  bool   `json:"success"`
		Response string `json:"response"`
	}{Success: true,
		Response: response}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respJson)
}
