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
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get job; Failed to grab auth token information: %v\n", err)
		return
	}

	// Grab search query from frontend (/api/jobs?search=QUERY)
	searchQuery := r.URL.Query().Get("search")

	jobs, err := db.GetJobs(tokenInfo.Uid, searchQuery)

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

// Handler for /api/jobs/{id} (GET)
func GetJob(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get job; Failed to grab auth token information: %v\n", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")

	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	job, err := db.GetJob(job_id, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get job: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)

}

// Handler for /api/jobs (PUT)
func UpdateJob(w http.ResponseWriter, r *http.Request) {
	var job db.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update job; Failed to grab decode request: %v\n", err)
		return
	}

	var tokenInfo Claim
	err, tokenInfo = GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update job; Failed to grab auth token information: %v\n", err)
		return
	}

	job.UserID = tokenInfo.Uid
	err = db.UpdateJob(job)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to update job: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully updated."}`)
}

// Handler for /api/jobs/{id}/archive (POST)
func ArchiveJob(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to archive job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to archive job; Failed to grab auth token information: %v\n", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")

	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	var job db.Job
	job.ID = job_id
	job.UserID = tokenInfo.Uid
	err = db.ArchiveJob(job)
	if err != nil {
		http.Error(w, "Failed to archive job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to archive job: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully Archived."}`)
}

// Handler for /api/jobs/{id}/unarchive (POST)
func UnarchiveJob(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to unarchive job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to unarchive job; Failed to grab auth token information: %v\n", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")

	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	var job db.Job
	job.ID = job_id
	job.UserID = tokenInfo.Uid
	err = db.UnarchiveJob(job)
	if err != nil {
		http.Error(w, "Failed to unarchive job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to unarchive job: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully Unarchived."}`)
}

// Handler for /api/jobs/{id} (Delete)
func DeleteJob(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to delete job; Failed to grab auth token information: %v\n", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")

	job_id, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to convert user id into integer: %v\n", err)
		return
	}

	var job db.Job
	job.ID = job_id
	job.UserID = tokenInfo.Uid
	err = db.DeleteJob(job)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to delete job: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully deleted."}`)
}

// Handler for /api/jobs/{id}/activities (GET)
func GetJobActivities(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Fprintf(os.Stderr, "Auth error: %v\n", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")

	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Invalid job id", http.StatusBadRequest)
		return
	}

	isJobOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil || !isJobOwner {
		http.Error(w, "Failed to verify ownership of job", http.StatusUnauthorized)
		return
	}

	activities, err := db.GetJobActivities(jobID)
	if err != nil {
		http.Error(w, "Failed to get job activities", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get job activities: %v\n", err)
		return
	}

	if activities == nil {
		activities = []db.JobActivity{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func SaveAIDocumentToJob(w http.ResponseWriter, r *http.Request) {
	err, token := GrabToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Invalid job id", http.StatusBadRequest)
		return
	}
	var body struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.Content == "" {
		http.Error(w, "empty content", http.StatusBadRequest)
		return
	}

	doc := db.Document{
		UserID:       token.Uid,
		Title:        fmt.Sprintf("AI %s - Job %d", body.Type, jobID),
		DocumentType: body.Type,
		IsArchived:   false,
	}

	docID, err := db.CreateDocument(doc)
	if err != nil {
		http.Error(w, "failed to create document", http.StatusInternalServerError)
		return
	}

	doc.ID = docID

	filePath := fmt.Sprintf("./data/documents/%d.txt", docID)

	err = os.WriteFile(filePath, []byte(body.Content), 0644)
	if err != nil {
		_ = db.DeleteDocument(token.Uid, int(docID))
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	_, err = db.CreateDocumentLink(jobID, docID, body.Type)
	if err != nil {
		http.Error(w, "Failed to link document", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"document": doc,
	})
}
