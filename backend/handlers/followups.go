package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
)

func CreateFollowUp(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		settings.Logger.Error("Failed to convert user id into integer", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to create follow up for job they don't own", "err", err)
		return
	}

	var input db.FollowUpTask
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		settings.Logger.Error("Failed to decode follow up creation request", "err", err)
		return
	}

	created, err := db.CreateFollowUp(jobID, input)
	if err != nil {
		http.Error(w, "Failed to create follow-up", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create follow up", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func GetFollowUps(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job follow-ups", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job follow-ups: Failed to convert job id into integer", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to get follow-ups for job they don't own", "err", err)
		return
	}

	tasks, err := db.GetFollowUps(jobID)
	if err != nil {
		http.Error(w, "Failed to get job follow-ups", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get job follow-ups", "err", err)
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
		http.Error(w, "Failed to update job follow-ups", http.StatusBadRequest)
		settings.Logger.Error("Failed to update job follow-ups: Failed to convert job id into integer", "err", err)
		return
	}

	follow_id_raw := chi.URLParam(r, "followup_id")
	followID, err := strconv.Atoi(follow_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		settings.Logger.Error("Failed to update job follow-up: Failed to convert follow-up id into integer", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to update follow up for job they don't own", "err", err)
		return
	}

	var input db.FollowUpTask
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		settings.Logger.Error("Failed to update job follow-up: Failed to decode follow up update request", "err", err)
		return
	}

	updated, err := db.UpdateFollowUp(followID, jobID, input)
	if err != nil {
		http.Error(w, "Failed to update follow-up", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update job follow-up", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func DeleteFollowUp(w http.ResponseWriter, r *http.Request) {
	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete job follow-up", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete job follow-up: Failed to convert job id into integer", "err", err)
		return
	}

	follow_id_raw := chi.URLParam(r, "followup_id")
	followID, err := strconv.Atoi(follow_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete job follow-up", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete job follow-up: Failed to convert follow-up id into integer", "err", err)
		return
	}
	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized Operation: User does not own job", http.StatusUnauthorized)
		settings.Logger.Info("Unauthorized Operation: User attempted to delete follow up for job they don't own", "err", err)
		return
	}

	err = db.DeleteFollowUp(followID, jobID)
	if err != nil {
		http.Error(w, "Failed to delete job follow-up", http.StatusInternalServerError)
		settings.Logger.Info("Failed to delete job follow-up", "err", err)
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
