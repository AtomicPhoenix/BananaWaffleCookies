package db

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Document struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Title        string    `json:"title"`
	DocumentType string    `json:"document_type"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func CreateDocument(doc Document) (int, error) {
	var id int
	sql_query := `INSERT INTO documents (user_id, title, document_type, is_archived) 
				VALUES ($1, $2, $3, $4) 
				RETURNING id`
	err := DbConn.QueryRow(context.Background(), sql_query, doc.UserID, doc.Title, doc.DocumentType, doc.IsArchived).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert document into database: %v\n", err)
		return -1, err
	}
	return id, err
}

func DeleteDocument(user_id int, doc_id int) error {
	sql_query := `DELETE FROM documents WHERE id = $1 AND user_id = $2`
	_, err := DbConn.Exec(context.Background(), sql_query, doc_id, user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove document from database: %v\n", err)
		return err
	}
	return nil
}

func UpdateDocument(doc Document) error {
	sql_query := `UPDATE documents
		      SET title = $1, document_type = $2, is_archived = $3
		      WHERE id = $4 AND user_id = $5`
	_, err := DbConn.Exec(context.Background(), sql_query, doc.Title, doc.DocumentType, doc.IsArchived, doc.ID, doc.UserID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update document in Database: %v\n", err)
		return err
	}
	return nil
}
