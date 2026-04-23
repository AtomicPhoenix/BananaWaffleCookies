package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
)

// POST /api/jobs/{id}/documents
func LinkDocumentToJob(w http.ResponseWriter, r *http.Request) {
	err, token := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to link document to job", http.StatusBadRequest)
		settings.Logger.Error("Failed to link document to job; Failed to grab auth token information", "err", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to link document to job", http.StatusBadRequest)
		settings.Logger.Error("Failed to link document to job; Failed to convert job id to int", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Failed to link document to job: Failed to verify job ownership", http.StatusUnauthorized)
		settings.Logger.Error("Failed to link document to job: Failed to verify job ownership", "err", err)
		return
	}

	var body struct {
		DocumentID int    `json:"document_id"`
		LinkType   string `json:"link_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to link document to job: Bad request body", http.StatusInternalServerError)
		settings.Logger.Error("Failed to link document to job: Bad request body", "err", err)
		return
	}

	_, err = db.GetDocument(int(body.DocumentID), token.Uid)
	if err != nil {
		http.Error(w, "Failed to link document to job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to link document to job: Could not get job", "err", err)
		return
	}

	id, err := db.CreateDocumentLink(jobID, body.DocumentID, body.LinkType)
	if err != nil {
		http.Error(w, "Failed to link document to job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to link document to job", "err", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"id":      id,
	})
}

// GET /api/jobs/{id}/documents
func GetJobDocuments(w http.ResponseWriter, r *http.Request) {
	err, token := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get job documents", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job documents; Failed to grab auth token information", "err", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job documents", http.StatusBadRequest)
		settings.Logger.Error("Failed to get job documents; Failed to convert job id to int", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Failed to get job documents: Failed to verify job ownership", http.StatusUnauthorized)
		settings.Logger.Error("Failed to get job documents: Failed to verify job ownership", "err", err)
		return
	}

	docs, err := db.GetJobDocuments(jobID, token.Uid)
	if err != nil {
		http.Error(w, "Failed to get job documents", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get job documents", "err", err)
		return
	}

	json.NewEncoder(w).Encode(docs)
}

// DELETE /api/jobs/{id}/documents/{document_id}
func UnlinkDocumentFromJob(w http.ResponseWriter, r *http.Request) {
	err, token := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to unlink document to job", http.StatusBadRequest)
		settings.Logger.Error("Failed to unlink document to job; Failed to grab auth token information", "err", err)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to unlink document to job", http.StatusBadRequest)
		settings.Logger.Error("Failed to unlink document to job; Failed to convert job id to int", "err", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Failed to unlink document to job: Failed to verify job ownership", http.StatusUnauthorized)
		settings.Logger.Error("Failed to unlink document to job: Failed to verify job ownership", "err", err)
		return
	}

	docID, err := strconv.Atoi(chi.URLParam(r, "document_id"))
	if err != nil {
		http.Error(w, "Failed to unlink document to job", http.StatusInternalServerError)
		settings.Logger.Error("Failed to unlink document to job: Failed to convert document id to int", "err", err)
		return
	}

	var job db.Job
	job.ID = jobID
	job.UserID = token.Uid

	var doc db.Document
	doc.ID = docID

	err = db.DeleteDocumentLink(job, doc)
	if err != nil {
		http.Error(w, "Failed to unlink document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to unlink document", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
