package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bananawafflecookies.com/m/v2/ai"
	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
)

// Handler for /api/jobs/{id}/resume (GET)
func GetResumeDraft(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusBadRequest)
		settings.Logger.Error("Failed to generate resume; Failed to grab auth token information", "err", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusBadRequest)
		settings.Logger.Error("Failed to generate resume; Failed to convert user id into integer", "err", err)
		return
	}

	job, err := db.GetJob(job_id, tokenInfo.Uid)

	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate resume; Failed to get job information", "err", err)
		return
	}

	profile, err := db.GetProfile(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate resume; Failed to get profile information", "err", err)
		return
	}

	response, err := ai.GenerateResumeDraft(job, profile)
	if err != nil {
		http.Error(w, "Failed to generate resume", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate resume", "err", err)
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

// Handler for /api/jobs/{id}/cover-letter (GET)
func GetCoverLetterDraft(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to generate cover letter", http.StatusBadRequest)
		settings.Logger.Error("Failed to generate cover letter. Failed to grab auth token", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		settings.Logger.Error("Failed to generate cover letter. Failed to convert job id into integer", "err", err)
		return
	}

	job, err := db.GetJob(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to generate cover letter", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate cover letter. Failed to get job info", "err", err)
		return
	}

	profile, err := db.GetProfile(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Internal server error when generating cover letter", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate cover letter. Failed to get profile info", "err", err)
		return
	}

	response, err := ai.GenerateCoverLetter(job, profile)
	if err != nil {
		http.Error(w, "Internal server error when generating cover letter", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate cover letter", "err", err)
		return
	}

	respJson := struct {
		Success  bool   `json:"success"`
		Response string `json:"response"`
	}{
		Success:  true,
		Response: response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respJson)
}
