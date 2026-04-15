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

// POST /api/jobs/{id}/documents
func LinkDocumentToJob(w http.ResponseWriter, r *http.Request) {
	err, token := GrabToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	var body struct {
		DocumentID int    `json:"document_id"`
		LinkType   string `json:"link_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = db.GetDocument(int(body.DocumentID), token.Uid)
	if err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}

	id, err := db.CreateDocumentLink(jobID, body.DocumentID, body.LinkType)
	if err != nil {
		http.Error(w, "failed to link document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "link error: %v\n", err)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	docs, err := db.GetJobDocuments(jobID, token.Uid)
	if err != nil {
		http.Error(w, "failed to fetch documents", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(docs)
}

// DELETE /api/jobs/{id}/documents/{document_id}
func UnlinkDocumentFromJob(w http.ResponseWriter, r *http.Request) {
	err, _ := GrabToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	job_id_raw := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(job_id_raw)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert job id into integer: %v\n", err)
		return
	}

	_, ok := verifyJobOwner(r, jobID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	docID, err := strconv.Atoi(chi.URLParam(r, "document_id"))
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to convert document id into integer: %v\n", err)
		return
	}

	err = db.DeleteDocumentLink(jobID, docID)
	if err != nil {
		http.Error(w, "Failed to unlink document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "unlink error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
