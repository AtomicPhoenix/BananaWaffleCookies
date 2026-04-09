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
