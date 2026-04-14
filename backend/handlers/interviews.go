package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"github.com/go-chi/chi/v5"
)

func CreateInterview(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to post job; Failed to grab auth token information: %v\n", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to verify ownership", http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var interview db.Interviews
	if err := json.NewDecoder(r.Body).Decode(&interview); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	interview.JobID = int(jobID)

	id, err := db.CreateInterview(interview)
	if err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
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
		http.Error(w, "Failed to post job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to post job; Failed to grab auth token information: %v\n", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Ownership check
	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to verify ownership", http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	interviews, err := db.GetInterviews(jobID)
	if err != nil {
		http.Error(w, "Failed to fetch interviews", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(interviews)
}

func DeleteInterview(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to post job; Failed to grab auth token information: %v\n", err)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	interviewIDStr := chi.URLParam(r, "interview_id")
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		http.Error(w, "Invalid interview ID", http.StatusBadRequest)
		return
	}

	// Ownership check
	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to verify ownership", http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = db.DeleteInterview(interviewID, jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete interview: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
