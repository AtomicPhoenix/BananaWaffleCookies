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

func CreateFollowUp(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input db.FollowUpTask
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	created, err := db.CreateFollowUp(jobID, input)
	if err != nil {
		http.Error(w, "Failed to create follow-up", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to create follow-up: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func GetFollowUps(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := db.GetFollowUps(jobID)
	if err != nil {
		http.Error(w, "Failed to fetch follow-ups", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to fetch follow-up: %v\n", err)
		return
	}

	if tasks == nil {
		tasks = []db.FollowUpTask{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func UpdateFollowUp(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	follow_id_raw := chi.URLParam(r, "followup_id")
	followID, err := strconv.Atoi(follow_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input db.FollowUpTask
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	updated, err := db.UpdateFollowUp(followID, jobID, input)
	if err != nil {
		http.Error(w, "Failed to update follow-up", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to update follow-up: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func DeleteFollowUp(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	follow_id_raw := chi.URLParam(r, "followup_id")
	followID, err := strconv.Atoi(follow_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}
	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = db.DeleteFollowUp(followID, jobID)
	if err != nil {
		http.Error(w, "Failed to delete follow-up", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func verifyJobOwner(r *http.Request, jobID int) (Claim, bool) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		return Claim{}, false
	}

	ok, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil || !ok {
		return Claim{}, false
	}

	return tokenInfo, true
}
