package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
)

func CreateInterview(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to create interview", http.StatusBadRequest)
		settings.Logger.Error("Failed to create interview; Failed to grab auth token information", "err", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Failed to create interview: Invalid job ID", http.StatusBadRequest)
		settings.Logger.Error("Failed to create interview; Failed to convert job id to int", "err", err)
		return
	}

	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to create interview: Failed to verify job ownership", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create interview: Failed to verify job ownership", "err", err)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to create interview for job they don't own", "err", err)
		return
	}

	var interview db.Interviews
	if err := json.NewDecoder(r.Body).Decode(&interview); err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create interview: Failed to decode job creation request", "err", err)
		return
	}

	interview.JobID = int(jobID)

	id, err := db.CreateInterview(interview)
	if err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create interview", "err", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"id": id,
	})
}

func GetInterviews(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get interviews", http.StatusBadRequest)
		settings.Logger.Error("Failed to get interviews; Failed to grab auth token information", "err", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Failed to get interviews", http.StatusBadRequest)
		settings.Logger.Error("Failed to get interviews; Failed to convert job id to int", "err", err)
		return
	}

	// Ownership check
	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get interviews", http.StatusBadRequest)
		settings.Logger.Error("Failed to get interviews; Failed to verify job ownership", "err", err)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to get interviews for job they don't own", "err", err)
		return
	}

	interviews, err := db.GetInterviews(jobID)
	if err != nil {
		http.Error(w, "Failed to get interviews", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get interviews", "err", err)
		return
	}

	json.NewEncoder(w).Encode(interviews)
}

func DeleteInterview(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to delete job interview", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete job interview; Failed to grab auth token information", "err", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Failed to delete interview", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete interviews; Failed to convert job id to int", "err", err)
		return
	}

	interviewIDStr := chi.URLParam(r, "interview_id")
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		http.Error(w, "Failed to delete interview", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete interview; Failed to convert interview id to int", "err", err)
		return
	}

	// Ownership check
	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to delete interview", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete interview; Failed to verify job ownership", "err", err)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to delete interview for job they don't own", "err", err)
		return
	}

	err = db.DeleteInterview(interviewID, jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to delete interview", http.StatusInternalServerError)
		settings.Logger.Error("Failed to delete interview", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
