package handlers

import (
	"context"
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
		settings.Logger.Error("Failed to upload document; Failed to copy file", "err", err)
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

// Handler for /api/documents/{id}/versions (POST)
func CreateDocumentVersion(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to create document version; Failed to grab auth token", "err", err)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Invalid document id", http.StatusBadRequest)
		settings.Logger.Error("Failed to create document version; Invalid document id", "err", err)
		return
	}

	if err := db.AssertDocumentOwnership(docID, tokenInfo.Uid); err != nil {
		http.Error(w, "Unauthroized", http.StatusForbidden)
		settings.Logger.Error("Failed to create document version; Ownership check failed", "err", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		settings.Logger.Error("Failed to create document version; Missing file", "err", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		settings.Logger.Error("Failed to create document version; Failed to read file", "err", err)
		return
	}

	if http.DetectContentType(buffer) != "application/pdf" {
		http.Error(w, "Only PDF allowed", http.StatusBadRequest)
		settings.Logger.Warn("Failed to create document version; Only PDFs allowed", "err", "User attempted to upload non-PDF file")
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to reset file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create document version; Failed to reset file reader", "err", err)
		return
	}

	filePath := fmt.Sprintf("./data/documents/%d.pdf", docID)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "File write failed", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create document version; Failed to create file", "err", err)
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		http.Error(w, "File save failed", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create document version; Failed to write file", "err", err)
		return
	}

	err = db.CreateDocumentVersion(docID, fileHeader.Filename, filePath, fileSize)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		settings.Logger.Error("Failed to create document version; DB insert failed", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for /api/documents/{id}/duplicate (POST)
func DuplicateDocument(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to duplicate document; Failed to grab auth token", "err", err)
		return
	}

	docIDRaw := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDRaw)
	if err != nil {
		http.Error(w, "Invalid document id", http.StatusBadRequest)
		settings.Logger.Error("Failed to duplicate document; Invalid document id", "err", err)
		return
	}

	tx, err := db.DbConn.Begin(context.Background())
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Failed to start transaction", "err", err)
		return
	}
	defer tx.Rollback(context.Background())

	// Grab current document
	var existing db.Document

	err = tx.QueryRow(
		context.Background(),
		`SELECT id, user_id, title, document_type, is_archived, current_version_id
		 FROM documents
		 WHERE id = $1 AND user_id = $2
		 FOR UPDATE`,
		docID,
		tokenInfo.Uid,
	).Scan(
		&existing.ID,
		&existing.UserID,
		&existing.Title,
		&existing.DocumentType,
		&existing.IsArchived,
		&existing.CurrentVersionID,
	)

	if err != nil {
		http.Error(w, "Document not found or not owned by user", http.StatusForbidden)
		settings.Logger.Error("Failed to duplicate document; Ownership or fetch failed", "err", err)
		return
	}

	// Grab latest document version
	var v db.DocumentVersion

	err = tx.QueryRow(
		context.Background(),
		`SELECT id, document_id, version_number, file_name, file_path, file_size_bytes
		 FROM document_versions
		 WHERE document_id = $1
		 ORDER BY version_number DESC
		 LIMIT 1`,
		docID,
	).Scan(
		&v.ID,
		&v.DocumentID,
		&v.VersionNumber,
		&v.FileName,
		&v.FilePath,
		&v.FileSizeBytes,
	)

	if err != nil {
		http.Error(w, "No version to duplicate", http.StatusBadRequest)
		settings.Logger.Error("Failed to duplicate document; No version found", "err", err)
		return
	}

	// Create new document row in database
	var newDocID int
	var newDocTitle string

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO documents (user_id, title, document_type, is_archived)
		 VALUES ($1, $2, $3, FALSE)
		 RETURNING id`,
		tokenInfo.Uid,
		newDocTitle,
		existing.DocumentType,
	).Scan(&newDocID)

	if err != nil {
		http.Error(w, "Failed to create duplicate document", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Insert document failed", "err", err)
		return
	}

	// Copy file on disk
	newFilePath := fmt.Sprintf("./data/documents/%d.pdf", newDocID)

	src, err := os.Open(v.FilePath)
	if err != nil {
		http.Error(w, "Failed to read source file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Failed to open source file", "err", err)
		return
	}
	defer src.Close()

	dst, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Failed to create destination file", "err", err)
		return
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, src)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; File copy failed", "err", err)
		return
	}

	// Create document version row in database for new document
	var newVersionID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO document_versions
		 (document_id, version_number, file_name, file_path, file_size_bytes)
		 VALUES (
			$1,
			1,
			$2,
			$3,
			$4
		 )
		 RETURNING id`,
		newDocID,
		v.FileName,
		newFilePath,
		fileSize,
	).Scan(&newVersionID)

	if err != nil {
		http.Error(w, "Failed to create version", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Insert version failed", "err", err)
		return
	}

	// Set current document version in database
	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		newVersionID,
		newDocID,
	)

	if err != nil {
		http.Error(w, "Failed to update document version", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Failed to update current version", "err", err)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		settings.Logger.Error("Failed to duplicate document; Transaction commit failed", "err", err)
		return
	}

	// Return new document
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":                 newDocID,
		"title":              newDocTitle,
		"current_version_id": newVersionID,
	})
}
