package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"github.com/go-chi/chi/v5"
)

// Handler for /api/documents (POST)
func UploadDocument(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to grab auth token information: %v\n", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to grab file from request: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Failed to parse document type", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to parse document type: %v\n", err)
		return
	}

	doc := db.Document{
		UserID:       tokenInfo.Uid,
		Title:        fileHeader.Filename,
		DocumentType: "resume",
		IsArchived:   false,
	}

	id, err := db.CreateDocument(doc)
	doc.ID = id

	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to create document entry in DB: %v\n", err)
		return
	}

	dst := fmt.Sprintf("./data/documents/%d.pdf", id)
	out, err := os.Create(dst)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to create file destination: %v\n", err)
		if err = db.DeleteDocument(tokenInfo.Uid, doc.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete document %s (ID: %d) from database: %v\n", doc.Title, doc.ID, err)
		}
		return
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to copy file to destination: %v\n", err)
		if err = db.DeleteDocument(tokenInfo.Uid, doc.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete document %s (ID: %d) from database: %v\n", doc.Title, doc.ID, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// Handler for /api/documents/{id} (DELETE)
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to delete document; Failed to grab auth token information: %v\n", err)
		return
	}

	doc_id_raw := chi.URLParam(r, "id")

	doc_id, err := strconv.Atoi(doc_id_raw)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to delete document: Failed to convert document id into integer: %v\n", err)
		return
	}
	err = db.DeleteDocument(tokenInfo.Uid, doc_id)
}

// Handler for /api/documents/{id} (PUT)
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to grab auth token information: %v\n", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to upload document; Failed to grab file from request: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Failed to parse document type", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to parse document type: %v\n", err)
		return
	}

	doc := db.Document{
		UserID:       tokenInfo.Uid,
		Title:        fileHeader.Filename,
		DocumentType: "resume",
		IsArchived:   false,
	}

	err = db.UpdateDocument(doc)
	if err != nil {
		http.Error(w, "Failed to update document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

// Handler for /api/documents/{id} (GET)
func GetDocument(w http.ResponseWriter, r *http.Request) {
	doc_id_raw := chi.URLParam(r, "id")

	doc_id, err := strconv.Atoi(doc_id_raw)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get document; Failed to convert user id into integer: %v\n", err)
		return
	}

	doc, err := db.GetDocument(doc_id)

	if err != nil {
		http.Error(w, "Failed to get document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get document: %v\n", err)
		return
	}

	// Grab document file
	filePath := fmt.Sprintf("./data/documents/%d.pdf", doc_id)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusNotFound)
		fmt.Fprintf(os.Stderr, "Failed to get document; Failed to open document file: %v\n", err)
		return
	}
	defer file.Close()

	// Get file's stats
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get document; Failed to stat file: %v\n", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.Title))
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	// Stream file to requesting user
	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send file: %v\n", err)
	}
}

// Handler for /api/documents/{id}/info (GET)
func GetDocumentInfo(w http.ResponseWriter, r *http.Request) {
	doc_id_raw := chi.URLParam(r, "id")

	doc_id, err := strconv.Atoi(doc_id_raw)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get document; Failed to convert user id into integer: %v\n", err)
		return
	}

	doc, err := db.GetDocument(doc_id)

	if err != nil {
		http.Error(w, "Failed to get document", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get document: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}
