package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"github.com/go-chi/chi/v5"
)

// Handler for /api/documents (POST)
func UploadDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		settings.Logger.Error("Failed to upload document; Failed to grab auth token information", "err", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		settings.Logger.Error("Failed to upload document; Failed to grab file from request", "err", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		http.Error(w, "Failed to upload document", http.StatusBadRequest)
		settings.Logger.Error("Failed to upload document; Failed to read file", "err", err)
		return
	}
	buf = buf[:n]

	if http.DetectContentType(buf) != "application/pdf" {
		http.Error(w, "Only PDF files allowed", http.StatusBadRequest)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to reset file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to upload document; Failed to reset file reader", "err", err)
		return
	}

	doc := db.Document{
		UserID:       tokenInfo.Uid,
		Title:        fileHeader.Filename,
		DocumentType: "resume",
		IsArchived:   false,
	}

	docID, err := db.CreateDocument(doc)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		settings.Logger.Error("Failed to upload document; Failed to create document", "err", err)
		return
	}

	filePath := fmt.Sprintf("./data/documents/%d.pdf", docID)

	out, err := os.Create(filePath)
	if err != nil {
		if err = db.DeleteDocument(tokenInfo.Uid, doc.ID); err != nil {
			settings.Logger.Error("Failed to upload document; Failed cleanup delete document", "err", err)
		}
		http.Error(w, "File creation failed", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		if err = db.DeleteDocument(tokenInfo.Uid, doc.ID); err != nil {
			settings.Logger.Error("Failed to upload document; Failed cleanup delete document", "err", err)
		}
		http.Error(w, "File write failed", http.StatusInternalServerError)
		return
	}

	err = db.CreateDocumentVersion(docID, fileHeader.Filename, filePath, fileSize)
	if err != nil {
		settings.Logger.Error("Failed to upload document; Failed to create version", "err", err)
		if err = db.DeleteDocument(tokenInfo.Uid, doc.ID); err != nil {
			settings.Logger.Error("Failed to upload document; Failed cleanup delete document", "err", err)
		}
		http.Error(w, "Version creation failed", http.StatusInternalServerError)
		return
	}

	doc.ID = docID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// Handler for /api/documents/{id} (DELETE)
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete document; Failed to grab auth token information", "err", err)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to delete document; Failed to convert document id to int", "err", err)
		return
	}

	if err := db.DeleteDocument(tokenInfo.Uid, docID); err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to delete document", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for /api/documents/{id} (PUT)
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Invalid document id", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		http.Error(w, "Failed to parse document type", http.StatusBadRequest)
		settings.Logger.Error("Failed to update document; Failed to parse document type", "err", err)
		return
	}

	fileType := http.DetectContentType(buffer)
	if fileType != "application/pdf" {
		http.Error(w, "Only PDF allowed", http.StatusBadRequest)
		return
	}

	// Reset file reader
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to reset document reader", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update document; Failed to reset document reader", "err", err)
		return
	}

	// Save file to disk
	filePath := fmt.Sprintf("./data/documents/%d.pdf", docID)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update document; Failed to create file", "err", err)
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update document; Failed to write file", "err", err)
		return
	}

	err = db.UpdateDocument(tokenInfo.Uid, docID, fileHeader.Filename, filePath, fileSize)

	if err != nil {
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update document", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for /api/documents/{id} (GET)
func GetDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusBadRequest)
		settings.Logger.Error("Failed to get document; Failed to grab auth token information", "err", err)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get document; Failed to convert document id to int", "err", err)
		return
	}

	doc, err := db.GetDocument(docID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusBadRequest)
		settings.Logger.Error("Failed to get document", "err", err)
		return
	}

	filePath := fmt.Sprintf("./data/documents/%d.pdf", docID)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusNotFound)
		settings.Logger.Error("Failed to get document; Failed to open document file", "err", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get document; Failed to stat file", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.Title))
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	if _, err = io.Copy(w, file); err != nil {
		settings.Logger.Error("Failed to get document; Failed to stream file", "err", err)
	}
}

// Handler for /api/documents (GET)
func GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		settings.Logger.Error("Failed to get documents; Failed to grab auth token information", "err", err)
		return
	}

	docs, err := db.GetAllDocuments(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get documents; DB query failed", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(docs); err != nil {
		settings.Logger.Error("Failed to encode documents response", "err", err)
	}
}

// Handler for /api/documents/{id}/info (GET)
func GetDocumentInfo(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get document info", http.StatusBadRequest)
		settings.Logger.Error("Failed to get document info; Failed to grab auth token information", "err", err)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Failed to get document info", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get document info; Failed to convert document id to int", "err", err)
		return
	}

	doc, err := db.GetDocument(docID, tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get document info", http.StatusBadRequest)
		settings.Logger.Error("Failed to get document info", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}
