package db

import (
	"context"
	"fmt"
	"time"
)

type Document struct {
	ID               int             `json:"id"`
	UserID           int             `json:"user_id"`
	Title            string          `json:"title"`
	DocumentType     string          `json:"document_type"`
	Tags             []string        `json:"tags"`
	IsArchived       bool            `json:"is_archived"`
	CurrentVersionID int             `json:"current_version_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Version          DocumentVersion `json:"versions"`
}

type DocumentVersion struct {
	ID            int       `json:"id"`
	DocumentID    int       `json:"document_id"`
	VersionNumber int       `json:"version_number"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	CreatedAt     time.Time `json:"created_at"`
}

func CreateDocument(doc Document) (int, error) {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(context.Background())

	var docID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO documents (user_id, title, document_type, is_archived)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		doc.UserID,
		doc.Title,
		doc.DocumentType,
		doc.IsArchived,
	).Scan(&docID)

	if err != nil {
		return -1, fmt.Errorf("Failed to insert document for user_id=%d title=%s; Creating document failed: %v\n", doc.UserID, doc.Title, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, fmt.Errorf("Failed to insert document for user_id=%d title=%s; Transaction Failed: %v\n", doc.UserID, doc.Title, err)
	}
	return docID, err
}

func CreateDocumentVersion(docID int, fileName, filePath string, fileSize int64) error {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO document_versions
		 (document_id, version_number, file_name, file_path, file_size_bytes)
		 VALUES (
			$1,
			COALESCE((SELECT MAX(version_number) + 1 FROM document_versions WHERE document_id = $1), 1),
			$2,
			$3,
			$4
		 )
		 RETURNING id`,
		docID,
		fileName,
		filePath,
		fileSize,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Insert version failed: %v\n",
			docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Updating current version failed: %v\n",
			docID, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Transaction failed: %v\n",
			docID, err)
	}

	return nil
}

func DeleteDocument(user_id int, doc_id int) error {
	sql_query := `DELETE FROM documents WHERE id = $1 AND user_id = $2`

	_, err := DbConn.Exec(context.Background(), sql_query, doc_id, user_id)
	if err != nil {
		return fmt.Errorf("Failed to delete document for user_id=%d document_id=%d: %v\n", user_id, doc_id, err)
	}

	return nil
}

func UpdateDocument(userID, docID int, fileName, filePath string, fileSize int64) error {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var tmp int

	err = tx.QueryRow(
		context.Background(),
		`SELECT id
		 FROM documents
		 WHERE id = $1 AND user_id = $2
		 FOR UPDATE`,
		docID,
		userID,
	).Scan(&tmp)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; User does not own document or document does not exist: %v\n",
			userID, docID, err)
	}

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO document_versions
		 (document_id, version_number, file_name, file_path, file_size_bytes)
		 VALUES (
			$1,
			COALESCE((SELECT MAX(version_number) + 1 FROM document_versions WHERE document_id = $1), 1),
			$2,
			$3,
			$4
		 )
		 RETURNING id`,
		docID,
		fileName,
		filePath,
		fileSize,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; Could not insert new version into DB: %v\n",
			userID, docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; Could not update current version: %v\n",
			userID, docID, err)
	}

	return tx.Commit(context.Background())
}

func SetDocumentVersion(userID, docID, versionNumber int) error {
	// Begin transaction
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`SELECT dv.id
		 FROM document_versions dv
		 JOIN documents d ON d.id = dv.document_id
		 WHERE dv.document_id = $1
		   AND dv.version_number = $2
		   AND d.user_id = $3`,
		docID,
		versionNumber,
		userID,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to set document version for user_id=%d document_id=%d; Could not get version id: %v\n", userID, docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to set document version for user_id=%d document_id=%d; Could not set version id: %v\n", userID, docID, err)
	}

	return tx.Commit(context.Background())
}

func GetDocument(docID int, userID int) (Document, error) {
	var doc Document

	err := DbConn.QueryRow(
		context.Background(),
		`SELECT 
			d.id, d.user_id, d.title, d.document_type, d.is_archived, d.created_at, d.updated_at,
			v.id, v.version_number, v.file_name, v.file_path, v.file_size_bytes, v.created_at
		 FROM documents d
		 JOIN document_versions v ON v.id = d.current_version_id
		 WHERE d.id = $1 AND d.user_id = $2`,
		docID,
		userID,
	).Scan(
		&doc.ID, &doc.UserID, &doc.Title, &doc.DocumentType, &doc.IsArchived, &doc.CreatedAt, &doc.UpdatedAt,
		&doc.Version.ID, &doc.Version.VersionNumber, &doc.Version.FileName, &doc.Version.FilePath, &doc.Version.FileSizeBytes, &doc.Version.CreatedAt,
	)

	if err != nil {
		return Document{}, fmt.Errorf("Failed to get document for user_id=%d, doc_id=%d; Query failed: %v\n", userID, docID, err)
	}

	return doc, nil
}

func GetAllDocuments(userID int) ([]Document, error) {

	rows, err := DbConn.Query(
		context.Background(),
		`SELECT 
			d.id, d.user_id, d.title, d.document_type, d.is_archived, d.created_at, d.updated_at,
			v.id, v.version_number, v.file_name, v.file_path, v.file_size_bytes, v.created_at
		 FROM documents d
		 JOIN document_versions v ON v.id = d.current_version_id
		 WHERE d.user_id = $1
		 ORDER BY d.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Query failed: %v\n", userID, err)
	}
	defer rows.Close()

	var results []Document

	for rows.Next() {
		var r Document

		err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.Title,
			&r.DocumentType,
			&r.IsArchived,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.Version.ID,
			&r.Version.VersionNumber,
			&r.Version.FileName,
			&r.Version.FilePath,
			&r.Version.FileSizeBytes,
			&r.Version.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while scanning row: %v\n", userID, err)
		}

		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while executing query / reading results: %v\n", userID, err)
	}

	return results, nil
}
