package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
	"github.com/jung-kurt/gofpdf"
)

// Handler for /api/jobs (POST)
func CreateJob(w http.ResponseWriter, r *http.Request) {
	var job db.Job

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to post job; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to post job", http.StatusBadRequest)
		settings.Logger.Error("Failed to post job; Failed to grab auth token information", "err", err)
		return
	}

	job.UserID = tokenInfo.Uid
	if _, err := db.CreateJob(job); err != nil {
		http.Error(w, "Failed to post job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to post job", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully uploaded."}`)
}

// Handler for /api/jobs (GET)
func GetJobs(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job; Failed to grab auth token information", "err", err)
		return
	}

	searchQuery := r.URL.Query().Get("search")

	jobs, err := db.GetJobs(tokenInfo.Uid, searchQuery)
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusBadRequest)
		settings.Logger.Error("Failed to get jobs", "err", err)
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
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get job; Failed to convert job id to int", "err", err)
		return
	}

	job, err := db.GetJob(jobID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get job", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// Handler for /api/jobs (PUT)
func UpdateJob(w http.ResponseWriter, r *http.Request) {
	var job db.Job

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to update job; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update job", http.StatusBadRequest)
		settings.Logger.Error("Failed to update job; Failed to grab auth token information", "err", err)
		return
	}

	job.UserID = tokenInfo.Uid
	if err := db.UpdateJob(job); err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update job", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully updated."}`)
}

// Handler for /api/jobs/{id}/archive (POST)
func ArchiveJob(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to archive job", http.StatusBadRequest)
		settings.Logger.Error("Failed to archive job; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Failed to archive job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to archive job; Failed to convert job id to int", "err", err)
		return
	}

	var job db.Job
	job.ID = jobID
	job.UserID = tokenInfo.Uid

	if err := db.ArchiveJob(job); err != nil {
		http.Error(w, "Failed to archive job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to archive job", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully archived."}`)
}

// Handler for /api/jobs/{id}/unarchive (POST)
func UnarchiveJob(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to unarchive job", http.StatusBadRequest)
		settings.Logger.Error("Failed to unarchive job; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Failed to unarchive job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to unarchive job; Failed to convert job id to int", "err", err)
		return
	}

	var job db.Job
	job.ID = jobID
	job.UserID = tokenInfo.Uid

	if err := db.UnarchiveJob(job); err != nil {
		http.Error(w, "Failed to unarchive job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to unarchive job", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Job successfully unarchived."}`)
}

// Handler for /api/jobs/{id} (Delete)
func DeleteJob(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete job; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to delete job; Failed to convert job id to int", "err", err)
		return
	}

	var job db.Job
	job.ID = jobID
	job.UserID = tokenInfo.Uid

	if err := db.DeleteJob(job); err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to delete job", "err", err)
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
		settings.Logger.Error("Failed to get job activities; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Invalid job id", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job activities; Invalid job id", "err", err)
		return
	}

	isOwner, err := db.IsJobOwner(jobID, tokenInfo.Uid)
	if err != nil || !isOwner {
		http.Error(w, "Failed to verify ownership of job", http.StatusUnauthorized)
		settings.Logger.Error("Failed to get job activities; Failed to verify job ownership", "err", err)
		return
	}

	activities, err := db.GetJobActivities(jobID)
	if err != nil {
		http.Error(w, "Failed to get job activities", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get job activities", "err", err)
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to save AI document; Failed to grab auth token information", "err", err)
		return
	}

	jobIDRaw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(jobIDRaw)
	if err != nil {
		http.Error(w, "Invalid job id", http.StatusBadRequest)
		settings.Logger.Error("Failed to save AI document; Failed to convert job id to int", "err", err)
		return
	}

	var body struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		settings.Logger.Error("Failed to save AI document; Failed to decode request body", "err", err)
		return
	}

	if body.Content == "" {
		http.Error(w, "Empty content", http.StatusBadRequest)
		settings.Logger.Error("Failed to save AI document; Empty content")
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
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to save AI document; Failed to create document", "err", err)
		return
	}

	versionNumber := 1
	filePath := buildFilePath(token.Uid, docID, versionNumber)

	pdfBytes, err := generatePDFBytes(body.Content)
	if err != nil {
		_ = db.DeleteDocument(token.Uid, int(docID))
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		settings.Logger.Error("Failed to generate PDF", "err", err)
		return
	}

	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
		_ = db.DeleteDocument(token.Uid, int(docID))
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to save AI document; Failed to write file", "err", err)
		return
	}

	err = db.CreateDocumentVersion(
		docID,
		fmt.Sprintf("ai-%s.txt", body.Type),
		filePath,
		int64(len(body.Content)),
	)

	if err != nil {
		_ = db.DeleteDocument(token.Uid, docID)
		http.Error(w, "Failed to link document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to save AI document; Failed to link document", "err", err)
		return
	}

	if _, err := db.CreateDocumentLink(jobID, docID, body.Type); err != nil {
		http.Error(w, "Failed to link document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to save AI document; Failed to link document", "err", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"id":      docID,
	})
}

func generatePDFBytes(content string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	for _, line := range strings.Split(content, "\n") {
		pdf.MultiCell(0, 5, line, "", "", false)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
