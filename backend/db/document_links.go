package db

import (
	"context"
	"fmt"
	"os"
	"time"
)

type DocumentLinks struct {
	ID         int       `json:"id"`
	JobID      int       `json:"job_id"`
	DocumentID int       `json:"document_id"`
	LinkType   string    `json:"link_type"`
	CreatedAt  time.Time `json:"created_at"`
}

func CreateDocumentLink(jobID, documentID int, linkType string) (int, error) {
	var id int

	sql := `
		INSERT INTO document_links (job_id, document_id, link_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := DbConn.QueryRow(context.Background(), sql, jobID, documentID, linkType).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateDocumentLink failed for job_id=%d document_id=%d link_type=%s: %v\n", jobID, documentID, linkType, err)
		return 0, err
	}

	return id, nil
}

func DeleteDocumentLink(job Job, doc Document) error {
	sql := `
		DELETE FROM document_links
		WHERE job_id = $1 AND document_id = $2
	`

	_, err := DbConn.Exec(context.Background(), sql, job.ID, doc.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteDocumentLink failed for job_id=%d document_id=%d: %v\n", job.ID, doc.ID, err)
		return err
	}

	return nil
}

func GetJobDocuments(jobID int, userID int) ([]Document, error) {
	sql := `
		SELECT d.id, d.user_id, d.title, d.document_type, d.is_archived, d.created_at, d.updated_at
		FROM documents d
		JOIN document_links dl ON dl.document_id = d.id
		WHERE dl.job_id = $1 AND d.user_id = $2
		ORDER BY d.created_at DESC
	`

	rows, err := DbConn.Query(context.Background(), sql, jobID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetJobDocuments query failed for job_id=%d user_id=%d: %w", jobID, userID, err)
	}
	defer rows.Close()

	var docs []Document

	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.Title,
			&d.DocumentType,
			&d.IsArchived,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetJobDocuments scan failed for job_id=%d user_id=%d: %w", jobID, userID, err)
		}
		docs = append(docs, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetJobDocuments row iteration failed for job_id=%d user_id=%d: %w", jobID, userID, err)
	}

	return docs, nil
}

